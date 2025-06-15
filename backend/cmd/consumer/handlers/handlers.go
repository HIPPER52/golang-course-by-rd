package handlers

import (
	"context"
	"course_project/cmd/consumer/handlers/message"
	"course_project/internal/clients"
	"course_project/internal/constants/consumer"
	"course_project/internal/services/logger"
	"fmt"
)

type Handler struct {
	messageHandler *message.Handler
}

func NewHandler(clnts *clients.Clients) *Handler {
	return &Handler{
		messageHandler: message.NewHandler(clnts),
	}
}

func (h *Handler) Handle(ctx context.Context, env consumer.Envelope) {
	switch env.Type {
	case consumer.TypeSaveMessage:
		h.messageHandler.HandleSave(ctx, env.Payload)
	default:
		logger.Error(ctx, fmt.Errorf("unknown message type: %s", env.Type))
	}
}
