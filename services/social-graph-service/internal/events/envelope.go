package events

import "time"

type Envelope[T AllowedPayloads] struct {
	EventType      string    `json:"event_type"`
	EventVersion   int       `json:"event_version"`
	OccurredAt     time.Time `json:"occurred_at"`
	Producer       string    `json:"producer"`
	IdempotencyKey string    `json:"idempotency_key"`
	Payload        T         `json:"payload"`
}
type AllowedPayloads interface {
	FriendAcceptedPayload | FriendRequestedPayload | FriendCanceledPayload | FriendRejectedPayload | FriendRemovedPayload
}
