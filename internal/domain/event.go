package domain

import (
	"encoding/json"
	"time"
)

type Event struct {
	EventID   int64           `json:"event_id"`
	EventType string          `json:"event_type"`
	Payload   json.RawMessage `json:"payload"`
	CreatedAt time.Time       `json:"created_at"`
}
