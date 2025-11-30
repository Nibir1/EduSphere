// server/api/recommendations.go

package api

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	db "github.com/nibir1/go-fiber-postgres-REST-boilerplate/db/sqlc"
	"github.com/nibir1/go-fiber-postgres-REST-boilerplate/token"
)

// Input for POST /api/recommendations
type createRecoReq struct {
	TranscriptID int64 `json:"transcript_id"`
}

type coursePick struct {
	CourseID  int64  `json:"course_id"`
	Match     int    `json:"match"`
	Rationale string `json:"rationale"`
}
type recoJSON struct {
	Summary string       `json:"summary"`
	Courses []coursePick `json:"courses"`
}

// Build a prompt: feed transcript text + course catalog and ask for JSON
func (s *Server) buildRecoPrompt(ctx context.Context, transcript db.Transcript, courses []db.Course) string {
	var b strings.Builder
	b.WriteString("You are an academic advisor AI.\n")
	b.WriteString("Given the student's transcript text and the course catalog, select 5 relevant courses.\n")
	b.WriteString("Return STRICT JSON with this schema:\n")
	b.WriteString(`{"summary": string, "courses":[{"course_id": number, "match": number (0-100), "rationale": string}]}` + "\n\n")
	b.WriteString("=== Transcript ===\n")
	if transcript.TextExtracted.Valid {
		b.WriteString(transcript.TextExtracted.String)
	} else {
		b.WriteString("(no text extracted)")
	}
	b.WriteString("\n\n=== Course Catalog ===\n")
	for _, c := range courses {
		fmt.Fprintf(&b, "- ID:%d | Code:%s | Name:%s | Lang:%s | Grading:%s | Org:%s\n",
			c.ID, c.Code, c.Name, coalesce(c.Language, "N/A"), coalesce(c.GradingScale, "N/A"), coalesce(c.Organiser, "N/A"))
		if c.LearningOutcomes.Valid && c.LearningOutcomes.String != "" {
			fmt.Fprintf(&b, "  Learning outcomes: %s\n", c.LearningOutcomes.String)
		}
		if c.Prerequisites.Valid && c.Prerequisites.String != "" {
			fmt.Fprintf(&b, "  Prereq: %s\n", c.Prerequisites.String)
		}
	}
	b.WriteString("\nReturn only JSON. No extra commentary.\n")
	return b.String()
}

// POST /api/recommendations
func (s *Server) createRecommendation(c *fiber.Ctx) error {
	payload, ok := c.Locals(authorizationPayloadKey).(*token.Payload)
	if !ok || payload == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(errorResponse(fmt.Errorf("unauthorized")))
	}
	var req createRecoReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
	}
	tr, err := s.store.GetTranscript(c.Context(), req.TranscriptID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(errorResponse(err))
	}
	if tr.UserUsername != payload.Username {
		return c.Status(fiber.StatusForbidden).JSON(errorResponse(fmt.Errorf("forbidden")))
	}

	// Load a bounded slice of catalog (you can page later)
	courses, err := s.store.ListCourses(c.Context(), 100)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}
	if len(courses) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse(fmt.Errorf("no courses in catalog")))
	}

	prompt := s.buildRecoPrompt(c.Context(), tr, courses)

	// call Ollama
	msgs := []ollamaMessage{
		{Role: "system", Content: "You are a helpful academic advisor."},
		{Role: "user", Content: prompt},
	}
	raw, err := callOllamaChat(context.Background(), s.config.OllamaBaseURL, s.config.OllamaModel, msgs, true)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}

	var parsed recoJSON
	if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
		// If the model sent non-strict JSON, try to salvage (simple trim)
		raw = strings.TrimSpace(raw)
		if i := strings.Index(raw, "{"); i >= 0 {
			raw = raw[i:]
		}
		if j := strings.LastIndex(raw, "}"); j >= 0 {
			raw = raw[:j+1]
		}
		_ = json.Unmarshal([]byte(raw), &parsed) // best-effort
	}

	// Store recommendation
	payloadJSON, _ := json.Marshal(parsed)
	reco, err := s.store.CreateRecommendation(c.Context(), db.CreateRecommendationParams{
		UserUsername: payload.Username,
		TranscriptID: db.Int64ToNull(req.TranscriptID),
		Summary:      sqlStringOrNull(parsed.Summary),
		Payload:      payloadJSON,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}

	return c.JSON(fiber.Map{
		"id":        reco.ID,
		"summary":   parsed.Summary,
		"courses":   parsed.Courses,
		"createdAt": reco.CreatedAt,
	})
}

// GET /api/recommendations
func (s *Server) listRecommendations(c *fiber.Ctx) error {
	payload, ok := c.Locals(authorizationPayloadKey).(*token.Payload)
	if !ok || payload == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(errorResponse(fmt.Errorf("unauthorized")))
	}
	items, err := s.store.ListRecommendations(c.Context(), payload.Username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}
	return c.JSON(items)
}

// GET /api/recommendations/:id
func (s *Server) getRecommendation(c *fiber.Ctx) error {
	payload, ok := c.Locals(authorizationPayloadKey).(*token.Payload)
	if !ok || payload == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(errorResponse(fmt.Errorf("unauthorized")))
	}
	id, err := parseIDParam(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
	}

	reco, err := s.store.GetRecommendation(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(errorResponse(err))
	}
	if reco.UserUsername != payload.Username {
		return c.Status(fiber.StatusForbidden).JSON(errorResponse(fmt.Errorf("forbidden")))
	}
	return c.JSON(reco)
}
