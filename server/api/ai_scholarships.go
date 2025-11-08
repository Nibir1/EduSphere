package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	db "github.com/nibir1/go-fiber-postgres-REST-boilerplate/db/sqlc"
	"github.com/nibir1/go-fiber-postgres-REST-boilerplate/token"
)

// ScholarshipReco represents one AI-generated scholarship recommendation
type ScholarshipReco struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Match       float64 `json:"match"`
	Link        string  `json:"link"`
}

// POST /api/scholarships/generate
func (s *Server) generateScholarships(c *fiber.Ctx) error {
	payload, ok := c.Locals(authorizationPayloadKey).(*token.Payload)
	if !ok || payload == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(errorResponse(fmt.Errorf("unauthorized")))
	}

	// 1️⃣ Get latest transcript
	transcripts, err := s.store.ListTranscripts(c.Context(), payload.Username)
	if err != nil || len(transcripts) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(errorResponse(fmt.Errorf("no transcripts found")))
	}
	fullTr, err := s.store.GetTranscript(c.Context(), transcripts[0].ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}
	txText := strings.TrimSpace(fullTr.TextExtracted.String)
	if txText == "" {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse(fmt.Errorf("transcript has no extracted text")))
	}

	// 2️⃣ Perform Brave web search
	query := "scholarships for international students studying computer science OR artificial intelligence"
	webResults, werr := fetchBraveResults(query, s.config.BraveAPIKey, s.config.BraveAPIURL, s.config.WebSearchMaxResults)

	log.Println("------------------------------------------------------------")
	log.Printf("[DEBUG] Brave Search Results (%d):", len(webResults))
	for i, w := range webResults {
		log.Printf("%d) %s -> %s", i+1, w.Title, w.URL)
	}
	log.Println("------------------------------------------------------------")

	if werr != nil {
		log.Printf("[WEB] Brave web search failed: %v", werr)
	}
	if len(webResults) == 0 {
		log.Println("[WEB] No Brave results found — continuing with transcript only.")
	}

	// 3️⃣ Build AI prompt
	var sb strings.Builder
	sb.WriteString(`
You are an academic scholarship advisor.
Your task: recommend scholarships — NOT university courses or classes.
Use the student's transcript only to understand their background (e.g. Software Engineering, AI, Data Science).
From the provided web search results, find the most relevant scholarships for this profile.

Return ONLY scholarships (no courses, no degrees, no projects).
Each result must include:
- title (scholarship name)
- description (what it offers or who it's for)
- match (number 0–100)
- link (URL to the scholarship)

Respond only in valid JSON array format.
`)

	sb.WriteString("Transcript:\n\"\"\"\n")
	sb.WriteString(txText)
	sb.WriteString("\n\"\"\"\n\n")

	if len(webResults) > 0 {
		sb.WriteString("Scholarship Web Results:\n")
		for _, w := range webResults {
			sb.WriteString(fmt.Sprintf("- %s\n  Link: %s\n  About: %s\n",
				truncate(w.Title, 100),
				w.URL,
				truncate(w.Snippet, 200),
			))
		}
	}

	// 4️⃣ Ollama messages (force valid JSON array)
	messages := []ollamaMessage{
		{
			Role: "system",
			Content: `
You are a strict JSON generator.
Always return a valid JSON array, even if there is only one scholarship.
Do not include markdown, code fences, explanations, or text outside the array.
Each element must have:
- "title": string
- "description": string
- "match": number (0–100)
- "link": string (valid URL)

Example of valid output:
[
  {
    "title": "Google DeepMind AI Master's Scholarship",
    "description": "Supports AI students at top universities.",
    "match": 95.3,
    "link": "https://www.iie.org/programs/google-deepmind-ai-masters-scholarships/"
  }
]
			`,
		},
		{
			Role:    "user",
			Content: sb.String(),
		},
	}

	// 5️⃣ Call Ollama inference
	resp, err := callOllamaChat(c.Context(), s.config.OllamaBaseURL, s.config.OllamaModel, messages, true)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse(fmt.Errorf("ollama failed: %v", err)))
	}

	log.Println("------------------------------------------------------------")
	log.Println("[DEBUG] Raw Scholarship AI Response:")
	log.Println(resp)
	log.Println("------------------------------------------------------------")

	// 6️⃣ Parse JSON response (smart parsing with nested fallback)
	var recs []ScholarshipReco

	// Try direct array first
	if err := json.Unmarshal([]byte(resp), &recs); err != nil {
		// Try nested "scholarships" key before logging fallback
		var wrapper struct {
			Scholarships []ScholarshipReco `json:"scholarships"`
		}
		if jerr := json.Unmarshal([]byte(resp), &wrapper); jerr == nil && len(wrapper.Scholarships) > 0 {
			recs = wrapper.Scholarships
		} else {
			log.Printf("[DEBUG] JSON parse failed, fallback to manual extract: %v", err)
			recs = extractScholarshipJSON(resp)
		}
	}

	// 7️⃣ Clean and sanitize
	filtered := make([]ScholarshipReco, 0, len(recs))
	for _, r := range recs {
		r.Title = strings.TrimSpace(r.Title)
		r.Description = strings.TrimSpace(r.Description)
		r.Link = strings.TrimSpace(r.Link)
		if r.Title == "" {
			continue
		}
		filtered = append(filtered, r)
	}
	recs = filtered

	// Sort by match score (desc)
	sort.Slice(recs, func(i, j int) bool { return recs[i].Match > recs[j].Match })

	// 8️⃣ Persist top results
	ctx := context.Background()
	toStore := recs
	if len(toStore) > 5 {
		toStore = toStore[:5]
	}
	for _, r := range toStore {
		_, err := s.store.CreateScholarship(ctx, db.CreateScholarshipParams{
			UserUsername: payload.Username,
			Title:        r.Title,
			Description:  sql.NullString{String: r.Description, Valid: r.Description != ""},
			MatchScore:   sql.NullFloat64{Float64: r.Match, Valid: r.Match > 0},
			Link:         sql.NullString{String: r.Link, Valid: r.Link != ""},
		})
		if err != nil {
			log.Printf("[DB] Save scholarship failed: %v", err)
		}
	}

	// 9️⃣ Log parsed results
	log.Println("------------------------------------------------------------")
	log.Println("[DEBUG] Parsed Scholarship Recommendations:")
	for i, r := range recs {
		log.Printf("%d) %s (%.2f%%) — %s", i+1, r.Title, r.Match, r.Link)
	}
	log.Println("------------------------------------------------------------")

	return c.JSON(fiber.Map{
		"user":         payload.Username,
		"scholarships": recs,
		"generated_at": time.Now(),
	})
}

// Helper: extract array if Ollama wraps JSON in text
func extractScholarshipJSON(raw string) []ScholarshipReco {
	start := strings.Index(raw, "[")
	end := strings.LastIndex(raw, "]")
	if start == -1 || end == -1 || end <= start {
		return nil
	}
	sub := raw[start : end+1]
	var recs []ScholarshipReco
	_ = json.Unmarshal([]byte(sub), &recs)
	return recs
}
