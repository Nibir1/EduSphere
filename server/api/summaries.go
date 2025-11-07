package api

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jung-kurt/gofpdf" // go get github.com/jung-kurt/gofpdf
	db "github.com/nibir1/go-fiber-postgres-REST-boilerplate/db/sqlc"
	"github.com/nibir1/go-fiber-postgres-REST-boilerplate/token"
)

type createSummaryReq struct {
	RecommendationID int64 `json:"recommendation_id"`
}

// POST /api/summaries
func (s *Server) createSummaryPDF(c *fiber.Ctx) error {
	payload, ok := c.Locals(authorizationPayloadKey).(*token.Payload)
	if !ok || payload == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(errorResponse(fmt.Errorf("unauthorized")))
	}
	var req createSummaryReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
	}

	reco, err := s.store.GetRecommendation(c.Context(), req.RecommendationID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(errorResponse(err))
	}
	if reco.UserUsername != payload.Username {
		return c.Status(fiber.StatusForbidden).JSON(errorResponse(fmt.Errorf("forbidden")))
	}

	// Generate PDF from reco.Summary + reco.Payload
	filename := fmt.Sprintf("summary_%d_%d.pdf", req.RecommendationID, time.Now().Unix())
	outPath := filepath.Join(s.summariesDir, filename)
	if err := writeRecoPDF(outPath, reco); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}

	row, err := s.store.CreateSummary(c.Context(), db.CreateSummaryParams{
		UserUsername:     payload.Username,
		RecommendationID: req.RecommendationID,
		PdfPath:          outPath,
	})
	if err != nil {
		_ = os.Remove(outPath)
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}

	return c.JSON(fiber.Map{
		"id":       row.ID,
		"pdf_path": row.PdfPath,
	})
}

// GET /api/summaries
func (s *Server) listSummaries(c *fiber.Ctx) error {
	payload, ok := c.Locals(authorizationPayloadKey).(*token.Payload)
	if !ok || payload == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(errorResponse(fmt.Errorf("unauthorized")))
	}
	items, err := s.store.ListSummaries(c.Context(), payload.Username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}
	return c.JSON(items)
}

// GET /api/summaries/:id/download
func (s *Server) downloadSummaryPDF(c *fiber.Ctx) error {
	payload, ok := c.Locals(authorizationPayloadKey).(*token.Payload)
	if !ok || payload == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(errorResponse(fmt.Errorf("unauthorized")))
	}
	id, err := parseIDParam(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
	}

	sum, err := s.store.GetSummary(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(errorResponse(err))
	}
	if sum.UserUsername != payload.Username {
		return c.Status(fiber.StatusForbidden).JSON(errorResponse(fmt.Errorf("forbidden")))
	}
	return c.Download(sum.PdfPath)
}

// DELETE /api/summaries/:id
func (s *Server) deleteSummary(c *fiber.Ctx) error {
	payload, ok := c.Locals(authorizationPayloadKey).(*token.Payload)
	if !ok || payload == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(errorResponse(fmt.Errorf("unauthorized")))
	}
	id, err := parseIDParam(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
	}

	sum, err := s.store.GetSummary(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(errorResponse(err))
	}
	if sum.UserUsername != payload.Username {
		return c.Status(fiber.StatusForbidden).JSON(errorResponse(fmt.Errorf("forbidden")))
	}

	_ = os.Remove(sum.PdfPath)
	if err := s.store.DeleteSummary(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// --- PDF writer ----
func writeRecoPDF(path string, reco db.Recommendation) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Helvetica", "B", 16)
	pdf.Cell(0, 10, "EduSphere: Recommended Courses Summary")
	pdf.Ln(12)

	pdf.SetFont("Helvetica", "", 12)
	if reco.Summary.Valid && reco.Summary.String != "" {
		pdf.MultiCell(0, 6, "Summary:\n"+reco.Summary.String, "", "", false)
		pdf.Ln(4)
	}
	pdf.MultiCell(0, 6, "Created at: "+reco.CreatedAt.Format(time.RFC3339), "", "", false)
	pdf.Ln(6)

	// dump top picks if present
	var payload struct {
		Summary string `json:"summary"`
		Courses []struct {
			CourseID  int64  `json:"course_id"`
			Match     int    `json:"match"`
			Rationale string `json:"rationale"`
		} `json:"courses"`
	}
	_ = json.Unmarshal(reco.Payload, &payload)

	if len(payload.Courses) > 0 {
		pdf.SetFont("Helvetica", "B", 14)
		pdf.Cell(0, 8, "Top Courses")
		pdf.Ln(10)
		pdf.SetFont("Helvetica", "", 11)
		for i, c := range payload.Courses {
			pdf.MultiCell(0, 6, fmt.Sprintf("%d) Course ID: %d | Match: %d%%\nRationale: %s",
				i+1, c.CourseID, c.Match, c.Rationale), "", "", false)
			pdf.Ln(2)
		}
	}
	return pdf.OutputFileAndClose(path)
}
