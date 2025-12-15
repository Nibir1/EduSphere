package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	db "github.com/nibir1/go-fiber-postgres-REST-boilerplate/db/sqlc"
	"github.com/nibir1/go-fiber-postgres-REST-boilerplate/token"
)

// Request payload from Frontend
type createRecommendationRequest struct {
	TranscriptID int64  `json:"transcript_id"`
	Preference   string `json:"preference"`
}

// Response struct
type Recommendation struct {
	Type        string  `json:"type"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Match       float64 `json:"match"`
	Code        string  `json:"code,omitempty"`
	Link        string  `json:"link,omitempty"`
	CourseID    int64   `json:"course_id,omitempty"`
}

// -----------------------------------------------------------------------------
// 1. HELPER: Extract Completed Courses (AI)
// -----------------------------------------------------------------------------
func (s *Server) extractCompletedCourses(c *fiber.Ctx, transcriptText string) ([]string, error) {
	messages := []aiMessage{
		{
			Role:    "system",
			Content: "You are a data extraction assistant. Analyze the academic transcript and return a JSON object with a single key 'completed_codes' containing a list of strings. Each string must be a Course Code (e.g. 'CS101') the student has completed.",
		},
		{
			Role:    "user",
			Content: transcriptText,
		},
	}

	raw, err := callOpenAIChat(c.Context(), s.config.OpenAIAPIKey, s.config.OpenAIModel, messages, true)
	if err != nil {
		return nil, err
	}

	var result struct {
		CompletedCodes []string `json:"completed_codes"`
	}
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		return []string{}, nil
	}

	normalized := make([]string, 0, len(result.CompletedCodes))
	for _, code := range result.CompletedCodes {
		normalized = append(normalized, strings.ToUpper(strings.TrimSpace(code)))
	}
	return normalized, nil
}

// -----------------------------------------------------------------------------
// 2. HELPER: Filter Logic
// -----------------------------------------------------------------------------
func filterAvailableCourses(allCourses []db.Course, completedCodes []string) []db.Course {
	completedMap := make(map[string]bool)
	for _, code := range completedCodes {
		completedMap[code] = true
	}

	var available []db.Course
	for _, course := range allCourses {
		dbCode := strings.ToUpper(strings.TrimSpace(course.Code))
		if !completedMap[dbCode] {
			available = append(available, course)
		}
	}
	return available
}

// -----------------------------------------------------------------------------
// 3. MAIN HANDLER: Create (Smart Filter)
// -----------------------------------------------------------------------------
func (s *Server) createRecommendation(c *fiber.Ctx) error {
	payload, ok := c.Locals(authorizationPayloadKey).(*token.Payload)
	if !ok || payload == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(errorResponse(fmt.Errorf("unauthorized")))
	}

	var req createRecommendationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
	}

	// Fetch Transcript
	transcript, err := s.store.GetTranscript(c.Context(), req.TranscriptID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(errorResponse(fmt.Errorf("transcript not found")))
	}
	if !transcript.TextExtracted.Valid || transcript.TextExtracted.String == "" {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse(fmt.Errorf("transcript has no text content")))
	}

	// AI Step A: Extract History
	completedCodes, err := s.extractCompletedCourses(c, transcript.TextExtracted.String)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse(fmt.Errorf("failed to analyze transcript history: %w", err)))
	}

	// DB: Get All Courses
	allCourses, err := s.store.ListAllCourses(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}

	// Go: Filter Available
	candidates := filterAvailableCourses(allCourses, completedCodes)
	if len(candidates) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"courses": []Recommendation{}, 
			"message": "No new courses available.",
		})
	}

	// Prepare AI Prompt
	type PromptCourse struct {
		ID   int64  `json:"id"`
		Code string `json:"code"`
		Name string `json:"name"`
		Desc string `json:"desc"`
	}
	var promptList []PromptCourse
	for _, c := range candidates {
		desc := ""
		if c.LearningOutcomes.Valid {
			desc = c.LearningOutcomes.String
			if len(desc) > 150 { desc = desc[:150] + "..." }
		}
		promptList = append(promptList, PromptCourse{
			ID:   c.ID,
			Code: c.Code,
			Name: c.Name,
			Desc: desc,
		})
	}
	candidateBytes, _ := json.Marshal(promptList)

	systemPrompt := `You are an academic course advisor.
	Task:
	1. Analyze the 'Available Courses' list and the 'User Preference'.
	2. Select the top 3-5 courses that best match the preference.
	3. Return a JSON object with a key "recommendations" which is an array.
	4. Each item must have: 
	   - "course_id" (integer, copied exactly from input)
	   - "code" (string)
	   - "title" (string)
	   - "rationale" (string, why it fits)
	   - "match" (number 0-100)`

	userPrompt := fmt.Sprintf("User Preference: %s\n\nAvailable Courses:\n%s", req.Preference, string(candidateBytes))

	messages := []aiMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}

	rawResponse, err := callOpenAIChat(c.Context(), s.config.OpenAIAPIKey, s.config.OpenAIModel, messages, true)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}

	// Parse AI Response
	var aiResult struct {
		Recommendations []struct {
			CourseID  int64   `json:"course_id"`
			Code      string  `json:"code"`
			Title     string  `json:"title"`
			Rationale string  `json:"rationale"`
			Match     float64 `json:"match"`
		} `json:"recommendations"`
	}
	
	if err := json.Unmarshal([]byte(rawResponse), &aiResult); err != nil {
		var directArr []struct {
			CourseID  int64   `json:"course_id"`
			Code      string  `json:"code"`
			Title     string  `json:"title"`
			Rationale string  `json:"rationale"`
			Match     float64 `json:"match"`
		}
		if err2 := json.Unmarshal([]byte(rawResponse), &directArr); err2 == nil {
			aiResult.Recommendations = directArr
		}
	}

	// Map back to Recommendation Struct
	finalRecs := make([]Recommendation, 0)
	for _, r := range aiResult.Recommendations {
		link := ""
		for _, c := range candidates {
			if c.ID == r.CourseID && c.CourseLink.Valid {
				link = c.CourseLink.String
				break
			}
		}

		finalRecs = append(finalRecs, Recommendation{
			Type:        "course",
			Title:       r.Title,
			Code:        r.Code,
			Description: r.Rationale,
			Match:       r.Match,
			Link:        link,
			CourseID:    r.CourseID,
		})
	}

	sort.Slice(finalRecs, func(i, j int) bool { return finalRecs[i].Match > finalRecs[j].Match })

	// Save to DB
	fullResultWrapper := fiber.Map{
		"courses": finalRecs,
	}
	resultJSON, _ := json.Marshal(fullResultWrapper)

	reco, err := s.store.CreateRecommendation(c.Context(), db.CreateRecommendationParams{
		UserUsername: payload.Username,
		TranscriptID: sql.NullInt64{Int64: req.TranscriptID, Valid: true},
		Payload:      resultJSON,
		Summary:      sql.NullString{String: "Course Recommendation", Valid: true},
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}

	return c.JSON(fiber.Map{
		"id":          reco.ID,
		"created_at":  reco.CreatedAt,
		"courses":     finalRecs,
		"user_pref":   req.Preference,
		"analyzed_at": time.Now(),
	})
}

// -----------------------------------------------------------------------------
// 4. HANDLERS: List & Get (Added these to fix build error)
// -----------------------------------------------------------------------------

// GET /recommendations
func (s *Server) listRecommendations(c *fiber.Ctx) error {
	payload, ok := c.Locals(authorizationPayloadKey).(*token.Payload)
	if !ok || payload == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(errorResponse(fmt.Errorf("unauthorized")))
	}

	// Calls the SQLC generated method ListRecommendations
	recos, err := s.store.ListRecommendations(c.Context(), payload.Username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}

	return c.JSON(recos)
}

// GET /recommendations/:id
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
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(errorResponse(fmt.Errorf("recommendation not found")))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}

	// Security check: ensure user owns this recommendation
	if reco.UserUsername != payload.Username {
		return c.Status(fiber.StatusForbidden).JSON(errorResponse(fmt.Errorf("forbidden")))
	}

	return c.JSON(reco)
}