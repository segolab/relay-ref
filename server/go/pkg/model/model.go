package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Destination struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

type CreateRelayRequest struct {
	EventType   string            `json:"eventType"`
	Destination Destination       `json:"destination"`
	Payload     json.RawMessage   `json:"payload"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

type RelayStatus string

const (
	RelayStatusQueued    RelayStatus = "queued"
	RelayStatusDelivered RelayStatus = "delivered"
	RelayStatusFailed    RelayStatus = "failed"
)

type Relay struct {
	ID            uuid.UUID         `json:"id"`
	EventType     string            `json:"eventType"`
	Destination   Destination       `json:"destination"`
	Payload       json.RawMessage   `json:"payload"`
	Metadata      map[string]string `json:"metadata,omitempty"`
	Status        RelayStatus       `json:"status"`
	CreatedAt     time.Time         `json:"createdAt"`
	DeliveredAt   *time.Time        `json:"deliveredAt,omitempty"`
	FailureReason *string           `json:"failureReason,omitempty"`
}

type ListRelaysResponse struct {
	Items         []*Relay `json:"items"`
	NextPageToken *string  `json:"nextPageToken"`
}

type ErrorResponse struct {
	Code      string         `json:"code"`
	Message   string         `json:"message"`
	Details   map[string]any `json:"details,omitempty"`
	RequestID *string        `json:"requestId,omitempty"`
}
