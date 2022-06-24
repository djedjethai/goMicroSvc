package service

import (
	"time"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type LogEntry struct {
	ID        string    `json:"id,omitempty"`
	Name      string    `json:"name"`
	Data      string    `json:"data"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
