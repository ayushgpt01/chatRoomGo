package chat

import "encoding/json"

type IncomingEvent struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type OutgoingEvent struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}
