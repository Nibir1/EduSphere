package api

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/nibir1/go-fiber-postgres-REST-boilerplate/token"
)

type chatReq struct {
	Messages []struct {
		Role    string `json:"role"`    // user|assistant|system
		Content string `json:"content"` // text
	} `json:"messages"`
}

func (s *Server) chatOnce(c *fiber.Ctx) error {
	_, ok := c.Locals(authorizationPayloadKey).(*token.Payload)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(errorResponse(fiber.ErrUnauthorized))
	}
	var req chatReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
	}
	if len(req.Messages) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse(fiber.ErrBadRequest))
	}

	// Map to Ollama schema
	msgs := make([]ollamaMessage, 0, len(req.Messages))
	for _, m := range req.Messages {
		msgs = append(msgs, ollamaMessage{Role: m.Role, Content: m.Content})
	}
	// Call model
	out, err := callOllamaChat(context.Background(), s.config.OllamaBaseURL, s.config.OllamaModel, msgs, false)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}
	return c.JSON(fiber.Map{"reply": out})
}
