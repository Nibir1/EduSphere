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

	// --- Build messages for Ollama ---
	ollamaMessages := []ollamaMessage{{
		Role: "system",
		Content: `You are EduSphere AI, an academic assistant.
Respond concisely and professionally, focusing on computer science, AI, and education.`,
	}}

	for _, m := range req.Messages {
		role := strings.ToLower(strings.TrimSpace(m.Role))
		if role != "user" && role != "assistant" {
			role = "user"
		}
		content := strings.TrimSpace(m.Content)
		if content != "" {
			ollamaMessages = append(ollamaMessages, ollamaMessage{
				Role:    role,
				Content: content,
			})
		}
	}

	// --- Setup streaming response headers ---
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")

	// Wrap writer so we can flush
	rawWriter := c.Response().BodyWriter()
	writer := bufio.NewWriter(rawWriter)

	// --- Prepare Ollama API request ---
	body := map[string]any{
		"model":    s.config.OllamaModel,
		"messages": ollamaMessages,
		"stream":   true,
	}
	bodyBytes, _ := json.Marshal(body)

	reqOllama, err := http.NewRequest("POST", s.config.OllamaBaseURL+"/api/chat", bytes.NewBuffer(bodyBytes))
	if err != nil {
		log.Printf("[CHAT-STREAM] Build request error: %v", err)
		fmt.Fprintf(writer, "event: error\ndata: build error\n\n")
		writer.Flush()
		return nil
	}
	reqOllama.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(reqOllama)
	if err != nil {
		log.Printf("[CHAT-STREAM] Ollama connection error: %v", err)
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

		if strings.HasPrefix(text, "{") {
			var data map[string]any
			if json.Unmarshal([]byte(text), &data) == nil {
				if msg, ok := data["message"].(map[string]any); ok {
					if token, ok := msg["content"].(string); ok && token != "" {
						fmt.Fprintf(writer, "data: %s\n\n", strings.ReplaceAll(token, "\n", "\\n"))
						writer.Flush() // ✅ works now
					}
				}
			}
		}
	}

	fmt.Fprint(writer, "data: [DONE]\n\n")
	writer.Flush()
	return nil
}
