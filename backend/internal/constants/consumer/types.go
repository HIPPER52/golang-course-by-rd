package consumer

import "encoding/json"

const TypeSaveMessage = "save_message"

type Envelope struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}
