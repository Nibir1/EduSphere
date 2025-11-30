// server/api/chat_ai.go

package api

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/nibir1/go-fiber-postgres-REST-boilerplate/token"
)

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// POST /api/chat/stream
func (s *Server) chatStream(c *fiber.Ctx) error {
	// --- Auth check ---
	payload, ok := c.Locals(authorizationPayloadKey).(*token.Payload)
	if !ok || payload == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(errorResponse(fmt.Errorf("unauthorized")))
	}

	// --- Parse request body ---
	var req struct {
		Messages []ChatMessage `json:"messages"`
	}
	if err := c.BodyParser(&req); err != nil || len(req.Messages) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse(fmt.Errorf("invalid request body")))
	}

	// --- Build messages for OpenAI ---
	var openAIMessages []map[string]string

	openAIMessages = append(openAIMessages, map[string]string{
		"role":    "system",
		"content": `You are EduSphere AI, an academic assistant. Respond concisely and professionally, focusing on computer science, AI, and education.`,
	})

	for _, m := range req.Messages {
		role := strings.ToLower(strings.TrimSpace(m.Role))
		if role != "user" && role != "assistant" && role != "system" {
			role = "user"
		}
		content := strings.TrimSpace(m.Content)
		if content != "" {
			openAIMessages = append(openAIMessages, map[string]string{
				"role":    role,
				"content": content,
			})
		}
	}

	// --- Setup streaming response headers ---
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")

	rawWriter := c.Response().BodyWriter()
	writer := bufio.NewWriter(rawWriter)

	// --- Prepare OpenAI streaming request ---
	body := map[string]any{
		"model":    s.config.OpenAIModel,
		"messages": openAIMessages,
		"stream":   true,
	}
	bodyBytes, _ := json.Marshal(body)

	reqOpenAI, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(bodyBytes))
	if err != nil {
		log.Printf("[CHAT-STREAM] Build OpenAI request error: %v", err)
		fmt.Fprintf(writer, "event: error\ndata: build error\n\n")
		writer.Flush()
		return nil
	}
	reqOpenAI.Header.Set("Content-Type", "application/json")
	reqOpenAI.Header.Set("Authorization", "Bearer "+s.config.OpenAIAPIKey)

	client := &http.Client{}
	resp, err := client.Do(reqOpenAI)
	if err != nil {
		log.Printf("[CHAT-STREAM] OpenAI connection error: %v", err)
		fmt.Fprintf(writer, "event: error\ndata: connection failed\n\n")
		writer.Flush()
		return nil
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err != io.EOF {
				log.Printf("[CHAT-STREAM] Read error: %v", err)
			}
			break
		}

		text := strings.TrimSpace(string(line))
		if text == "" {
			continue
		}

		if !strings.HasPrefix(text, "data: ") {
			continue
		}

		data := strings.TrimSpace(strings.TrimPrefix(text, "data: "))
		if data == "" {
			continue
		}

		if data == "[DONE]" {
			break
		}

		// Parse OpenAI streaming chunk
		var chunk struct {
			Choices []struct {
				Delta struct {
					Content string `json:"content"`
				} `json:"delta"`
			} `json:"choices"`
		}
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue
		}
		if len(chunk.Choices) == 0 {
			continue
		}

		token := chunk.Choices[0].Delta.Content
		if token == "" {
			continue
		}

		escaped := strings.ReplaceAll(token, "\n", "\\n")
		fmt.Fprintf(writer, "data: %s\n\n", escaped)
		writer.Flush()
	}

	fmt.Fprint(writer, "data: [DONE]\n\n")
	writer.Flush()
	return nil
}
