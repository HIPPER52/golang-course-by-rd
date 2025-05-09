package commands

import "encoding/json"

type CommandMessage struct {
	Command string          `json:"command"`
	Payload json.RawMessage `json:"payload"`
}
