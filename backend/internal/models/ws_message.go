package models

import "encoding/json"

type WSMessage struct {
	Event   string          `json:"event"`
	RoomID  string          `json:"room_id"`
	Payload json.RawMessage `json:"payload"`
}

type MessagePayload struct {
	Text string `json:"text"`
}

type IncomingMessage struct {
	Event   string `json:"event"`
	RoomID  string `json:"room_id"`
	Message string `json:"message"`
}
