package models

import (
	"time"
)

type DialogBase struct {
	ID            string    `json:"id" bson:"id"`
	ClientID      string    `json:"client_id" bson:"client_id"`
	ClientName    string    `json:"client_name" bson:"client_name"`
	ClientPhone   string    `json:"client_phone,omitempty" bson:"client_phone,omitempty"`
	ClientIP      string    `json:"client_ip,omitempty" bson:"client_ip,omitempty"`
	OperatorID    string    `json:"operator_id,omitempty" bson:"operator_id,omitempty"`
	StartedAt     time.Time `json:"started_at" bson:"started_at"`
	LastMessageAt time.Time `json:"last_message_at" bson:"last_message_at"`
	EndedAt       time.Time `json:"ended_at,omitempty" bson:"ended_at,omitempty"`
}
