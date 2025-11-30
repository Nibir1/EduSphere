// server/api/ai.go

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// Ollama message structure for chat requests
type ollamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Request/response structure
type ollamaChatReq struct {
	Model    string          `json:"model"`
	Messages []ollamaMessage `json:"messages"`
	Stream   bool            `json:"stream"`
	Format   string          `json:"format,omitempty"` // "json" for structured output
	Options  map[string]any  `json:"options,omitempty"`
}

type ollamaChatResp struct {
	Message ollamaMessage `json:"message"`
	Error   string        `json:"error,omitempty"`
}

// ------------------------------------------------------------------
// callOllamaChat: safer, longer timeout + better error diagnostics
// ------------------------------------------------------------------
func callOllamaChat(ctx context.Context, baseURL, model string, messages []ollamaMessage, expectJSON bool) (string, error) {
	url := fmt.Sprintf("%s/api/chat", baseURL)
	reqBody := ollamaChatReq{
		Model:    model,
		Messages: messages,
		Stream:   false,
		Options: map[string]any{
			"num_ctx":     4096, // larger context window
			"num_predict": 512,  // limit generation
			"temperature": 0.4,  // more deterministic
		},
	}

	if expectJSON {
		reqBody.Format = "json"
	}

	b, _ := json.Marshal(reqBody)

	httpClient := &http.Client{
		Timeout: 480 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(b))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	log.Printf("[AI] Sending request to Ollama: model=%s, len(messages)=%d, url=%s", model, len(messages), url)

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to reach Ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 500))
		return "", fmt.Errorf("ollama returned %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var out struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
		Error string `json:"error,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 500))
		log.Printf("[AI ERROR] failed to decode response: %v\nResponse body (truncated): %s", err, string(body))
		return "", fmt.Errorf("invalid Ollama response: %w", err)
	}

	if out.Error != "" {
		return "", fmt.Errorf("ollama error: %s", out.Error)
	}

	log.Printf("[AI] Ollama response (first 200 chars): %s", truncate(out.Message.Content, 200))
	return out.Message.Content, nil
}

// Helper: safe truncation for log output
func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "...(truncated)"
}
