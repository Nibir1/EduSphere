// server/api/ai_recommendations.go

package api

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nibir1/go-fiber-postgres-REST-boilerplate/token"
)

type Recommendation struct {
	Type        string  `json:"type"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Match       float64 `json:"match"`
}

// POST /api/recommendations/generate
func (s *Server) generateRecommendations(c *fiber.Ctx) error {
	payload, ok := c.Locals(authorizationPayloadKey).(*token.Payload)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(errorResponse(fmt.Errorf("unauthorized")))
	}

	// Get user's latest transcript
	transcripts, err := s.store.ListTranscripts(c.Context(), payload.Username)
	if err != nil || len(transcripts) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(errorResponse(fmt.Errorf("no transcripts found")))
	}

	latest := transcripts[0]
	fullTr, err := s.store.GetTranscript(c.Context(), latest.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}

	transcriptText := strings.TrimSpace(fullTr.TextExtracted.String)
	if transcriptText == "" {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse(fmt.Errorf("transcript has no extracted text")))
	}

	// Get all courses
	courses, err := s.store.ListAllCourses(context.Background())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}
	if len(courses) == 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse(fmt.Errorf("no courses in catalog")))
	}

	// Build course list context
	var courseList strings.Builder
	for _, cCourse := range courses {
		courseList.WriteString(fmt.Sprintf(
			"- %s (%s): %s\n",
			cCourse.Name,
			cCourse.Code,
			cCourse.LearningOutcomes.String,
		))
	}

	// System + user messages for OpenAI
	systemMsg := aiMessage{
		Role: "system",
		Content: `You are an academic advisor AI.
Given a student's transcript and a catalog of courses, recommend the most relevant university courses.

You MUST respond in JSON only. No markdown, no commentary.`,
	}

	userPrompt := fmt.Sprintf(`
The following text is an academic transcript. Analyze the student's background and recommend the most relevant university courses.

Transcript:
"""
%s
"""

Available courses:
%s

Return a JSON array of recommended courses like:
[
  {"title": "Course Name", "description": "Why it fits", "match": 95.2}
]
`, transcriptText, courseList.String())

	userMsg := aiMessage{
		Role:    "user",
		Content: userPrompt,
	}

	raw, err := callOpenAIChat(c.Context(), s.config.OpenAIAPIKey, s.config.OpenAIModel, []aiMessage{systemMsg, userMsg}, true)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse(fmt.Errorf("openai request failed: %v", err)))
	}

	// Parse model output robustly — support both array and object-shaped responses
	var recs []Recommendation

	// Try direct array first
	if err := json.Unmarshal([]byte(raw), &recs); err != nil {
		// If it's not an array, maybe it's an object with "courses" or "recommendations"
		var obj map[string]any
		if err2 := json.Unmarshal([]byte(raw), &obj); err2 == nil {
			if arr, ok := obj["courses"].([]any); ok {
				b, _ := json.Marshal(arr)
				_ = json.Unmarshal(b, &recs)
			} else if arr, ok := obj["recommendations"].([]any); ok {
				b, _ := json.Marshal(arr)
				_ = json.Unmarshal(b, &recs)
			}
		}

		// Fallback manual extraction if still empty
		if len(recs) == 0 {
			recs = extractJSONRecommendations(raw)
		}
	}

	// 🪶 DEBUG LOG: print parsed recommendation structs
	fmt.Println("[DEBUG] Parsed Recommendations (Go structs):")
	for i, r := range recs {
		fmt.Printf("  %d) %s (%.2f%%) — %s\n", i+1, r.Title, r.Match, r.Description)
	}
	fmt.Println("------------------------------------------------------------")

	// Sort by match score
	sort.Slice(recs, func(i, j int) bool { return recs[i].Match > recs[j].Match })

	// Split into categories (simple rule-based classification)
	var recommendedCourses []Recommendation
	var scholarships []Recommendation

	for _, r := range recs {
		if strings.Contains(strings.ToLower(r.Title), "scholarship") {
			r.Type = "scholarship"
			scholarships = append(scholarships, r)
		} else {
			r.Type = "course"
			recommendedCourses = append(recommendedCourses, r)
		}
	}

	// Return structured response
	return c.JSON(fiber.Map{
		"user":          payload.Username,
		"courses":       recommendedCourses,
		"scholarships":  scholarships,
		"analyzed_at":   time.Now(),
		"source":        s.config.OpenAIModel,
		"transcript_id": latest.ID,
	})
}

// Basic JSON parser fallback
func extractJSONRecommendations(raw string) []Recommendation {
	start := strings.Index(raw, "[")
	end := strings.LastIndex(raw, "]")
	if start == -1 || end == -1 || end <= start {
		return nil
	}
	sub := raw[start : end+1]
	var recs []Recommendation
	_ = json.Unmarshal([]byte(sub), &recs)
	return recs
}
