package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

type ollamaResponse struct {
	Response string `json:"response"`
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

	// Prepare prompt for Ollama
	var courseList strings.Builder
	for _, c := range courses {
		courseList.WriteString(fmt.Sprintf("- %s (%s): %s\n", c.Name, c.Code, c.LearningOutcomes.String))
	}

	prompt := fmt.Sprintf(`
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

	// Call Ollama API
	reqBody, _ := json.Marshal(map[string]any{
		"model":  "gemma3:4b-it-qat",
		"prompt": prompt,
		"stream": false,
	})
	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse(fmt.Errorf("ollama request failed: %v", err)))
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)

	// 🪶 DEBUG LOG: print the raw Ollama JSON response (first 500 chars)
	fmt.Println("------------------------------------------------------------")
	fmt.Println("[DEBUG] Raw Ollama Response:")
	fmt.Println(string(data))
	fmt.Println("------------------------------------------------------------")

	var parsed ollamaResponse
	if err := json.Unmarshal(data, &parsed); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse(fmt.Errorf("bad ollama response: %v", err)))
	}

	// Parse model output robustly — support both array and object-shaped responses
	var recs []Recommendation

	// Try direct array first
	if err := json.Unmarshal([]byte(parsed.Response), &recs); err != nil {
		// If it's not an array, maybe it's an object with "courses" or "recommendations"
		var obj map[string]any
		if err2 := json.Unmarshal([]byte(parsed.Response), &obj); err2 == nil {
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
			recs = extractJSONRecommendations(parsed.Response)
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
		"source":        "gemma3:4b-it-qat",
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
