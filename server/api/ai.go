package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Ollama request/response (chat)
type ollamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type ollamaChatReq struct {
	Model    string          `json:"model"`
	Messages []ollamaMessage `json:"messages"`
	Stream   bool            `json:"stream"`
	Format   string          `json:"format,omitempty"` // "json" if we want structured output
	Options  map[string]any  `json:"options,omitempty"`
}
type ollamaChatResp struct {
	Message ollamaMessage `json:"message"`
	// (other fields omitted)
}

func callOllamaChat(ctx context.Context, baseURL, model string, messages []ollamaMessage, expectJSON bool) (string, error) {
	url := fmt.Sprintf("%s/api/chat", baseURL)
	reqBody := ollamaChatReq{
		Model:    model,
		Messages: messages,
		Stream:   false,
	}
	if expectJSON {
		reqBody.Format = "json"
	}
	b, _ := json.Marshal(reqBody)

	httpClient := &http.Client{Timeout: 60 * time.Second}
	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("ollama returned %d", resp.StatusCode)
	}
	var out struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	return out.Message.Content, nil
}
