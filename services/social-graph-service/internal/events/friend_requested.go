package events

import (
	"encoding/json"
	"social-graph-service/internal/domain"
	"time"

	"github.com/google/uuid"
)

type Envelope[T AllowedPayloads] struct {
	EventType      string    `json:"event_type"`
	EventVersion   int       `json:"event_version"`
	OccurredAt     time.Time `json:"occurred_at"`
	Producer       string    `json:"producer"`
	IdempotencyKey string    `json:"idempotency_key"`
	Payload        T         `json:"payload"`
}
type AllowedPayloads interface {
	FriendAcceptedPayload | FriendRequestedPayload
}
type FriendRequestedPayload struct {
	RequestID  string `json:"request_id"`
	FromUserID string `json:"from_user_id"`
	ToUserID   string `json:"to_user_id"`
}

func BuildFriendRequested(reqID, fromUserID, toUserID uuid.UUID) (domain.OutboxEvent, error) {
	idempotencyKey := "friend.requested:" + reqID.String()
	envelope := Envelope[FriendRequestedPayload]{
		EventType:      "friend.requested",
		EventVersion:   1,
		OccurredAt:     time.Now().UTC(),
		Producer:       "social-graph-service",
		IdempotencyKey: idempotencyKey,
		Payload: FriendRequestedPayload{
			RequestID:  reqID.String(),
			FromUserID: fromUserID.String(),
			ToUserID:   toUserID.String()}}
	payload, err := json.Marshal(envelope)
	if err != nil {
		return domain.OutboxEvent{}, err
	}
	return domain.OutboxEvent{Topic: "socset.friends.events", PartitionKey: fromUserID.String(), Payload: payload}, nil
}
