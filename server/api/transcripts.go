package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	db "github.com/nibir1/go-fiber-postgres-REST-boilerplate/db/sqlc"
	"github.com/nibir1/go-fiber-postgres-REST-boilerplate/token"
)

// -----------------------------------------------------------------------------
// PDF TEXT EXTRACTION (stubs you can wire to ledongthuc/pdf or gosseract)
// -----------------------------------------------------------------------------

// Try plain text extraction using e.g. github.com/ledongthuc/pdf.
// Currently returns a sentinel error so the handler can optionally try OCR.
func extractPDFText(path string) (string, error) {
	// Example wiring (not included to keep deps light):
	// f, r, err := pdf.Open(path)
	// if err != nil { return "", err }
	// defer f.Close()
	// var buf strings.Builder
	// b, err := r.GetPlainText()
	// if err != nil { return "", err }
	// _, _ = io.Copy(&buf, b)
	// return buf.String(), nil
	return "", errors.New("text extractor not wired (plug ledongthuc/pdf or similar)")
}

// Optional OCR (requires system tesseract + gosseract binding).
// Currently returns a sentinel error unless you wire it up.
func ocrPDFToText(path string) (string, error) {
	return "", errors.New("OCR not implemented; wire gosseract if needed")
}

// -----------------------------------------------------------------------------
// HANDLERS
// -----------------------------------------------------------------------------

// POST /api/transcripts/upload  (multipart/form-data: file=<pdf>)
func (s *Server) uploadTranscript(c *fiber.Ctx) error {
	// 0) Auth
	payload, ok := c.Locals(authorizationPayloadKey).(*token.Payload)
	if !ok || payload == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(errorResponse(fmt.Errorf("unauthorized")))
	}

	// 1) Validate file presence and type
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse(fmt.Errorf("missing file: %w", err)))
	}
	if !strings.HasSuffix(strings.ToLower(fileHeader.Filename), ".pdf") {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse(fmt.Errorf("only PDF files are allowed")))
	}

	// 2) Save the file
	path, err := s.saveUploadedFile(fileHeader)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}

	// 3) Extract text (try normal first; optionally OCR if enabled)
	text, err := extractPDFText(path)
	meta := map[string]any{
		"ocr_used": false,
		"source":   "upload",
	}
	if err != nil || strings.TrimSpace(text) == "" {
		// OCR fallback if configured
		if s.config.OCRFallbackEnabled {
			if txt, oerr := ocrPDFToText(path); oerr == nil && strings.TrimSpace(txt) != "" {
				text = txt
				meta["ocr_used"] = true
			}
		}
	}

	// 4) Prepare DB payload
	metaJSON, _ := json.Marshal(meta)

	created, err := s.store.CreateTranscript(c.Context(), db.CreateTranscriptParams{
		UserUsername:  payload.Username,
		FilePath:      path,
		TextExtracted: sqlStringOrNull(text),
		Meta:          metaJSON,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":         created.ID,
		"file_path":  created.FilePath,
		"created_at": created.CreatedAt,
		"text_bytes": len(text),
		"ocr_used":   meta["ocr_used"],
	})
}

// GET /api/transcripts
func (s *Server) listTranscripts(c *fiber.Ctx) error {
	// 0) Auth
	payload, ok := c.Locals(authorizationPayloadKey).(*token.Payload)
	if !ok || payload == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(errorResponse(fmt.Errorf("unauthorized")))
	}

	items, err := s.store.ListTranscripts(c.Context(), payload.Username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}
	return c.JSON(items)
}

// GET /api/transcripts/:id
func (s *Server) getTranscript(c *fiber.Ctx) error {
	// 0) Auth
	payload, ok := c.Locals(authorizationPayloadKey).(*token.Payload)
	if !ok || payload == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(errorResponse(fmt.Errorf("unauthorized")))
	}

	// 1) Parse path param
	id, err := parseIDParam(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
	}

	// 2) Fetch transcript
	tr, err := s.store.GetTranscript(c.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(errorResponse(fmt.Errorf("not found")))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}

	// 3) Ownership check
	if tr.UserUsername != payload.Username {
		return c.Status(fiber.StatusForbidden).JSON(errorResponse(fmt.Errorf("forbidden")))
	}

	// 4) Prepare a short preview without slicing a non-string type
	preview := "(no text extracted)"
	if tr.TextExtracted.Valid {
		preview = tr.TextExtracted.String
		if len(preview) > 2000 {
			preview = preview[:2000] + "...(truncated)"
		}
	}

	return c.JSON(fiber.Map{
		"id":           tr.ID,
		"file_path":    tr.FilePath,
		"created_at":   tr.CreatedAt,
		"text_preview": preview,
	})
}
