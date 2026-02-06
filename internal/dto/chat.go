package dto

import "encoding/json"

type IncomingMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type OutgoingEvent struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}
