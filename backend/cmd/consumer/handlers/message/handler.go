package message

import (
	"context"
	"course_project/internal/clients"
	"course_project/internal/constants"
	"course_project/internal/models"
	"course_project/internal/services/logger"
	"encoding/json"
	"fmt"
	"time"
)

type Handler struct {
	clnts *clients.Clients
}

func NewHandler(clnts *clients.Clients) *Handler {
	return &Handler{clnts: clnts}
}

func (h *Handler) HandleSave(ctx context.Context, payload json.RawMessage) {
	var msg models.Message
	if err := json.Unmarshal(payload, &msg); err != nil {
		logger.Error(ctx, fmt.Errorf("failed to parse message payload: %w", err))
		return
	}

	if msg.SentAt.IsZero() {
		msg.SentAt = time.Now().UTC()
	}

	_, err := h.clnts.Mongo.Db.Collection(constants.CollectionMessages).InsertOne(ctx, msg)
	if err != nil {
		logger.Error(ctx, fmt.Errorf("failed to save message to MongoDB: %w", err))
		return
	}

	logger.Info(ctx, fmt.Sprintf("Message saved: %s", msg.ID))
}
