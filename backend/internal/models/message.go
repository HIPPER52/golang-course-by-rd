package models

import (
	"course_project/internal/constants/message"
	"time"
)

type Message struct {
	ID       string              `json:"id" bson:"id"`
	RoomID   string              `json:"room_id" bson:"room_id"`
	SenderID string              `json:"sender_id" bson:"sender_id"`
	Type     message.TypeMessage `json:"type" bson:"type"`
	Content  string              `json:"content" bson:"content"`
	SentAt   time.Time           `json:"sent_at" bson:"sent_at"`
}
