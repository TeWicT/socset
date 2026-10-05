package events

import (
	"encoding/json"
	"social-graph-service/internal/domain"
	"time"

	"github.com/google/uuid"
)

type FriendAcceptedPayload struct {
	RequestID  string `json:"request_id"`
	FromUserID string `json:"from_user_id"`
	ToUserID   string `json:"to_user_id"`
}

func BuildFriendAccepted(reqID, fromUserID, toUserID uuid.UUID) (domain.OutboxEvent, error) {
	idempotencyKey := "friend.accepted:" + reqID.String()
	envelope := Envelope[FriendAcceptedPayload]{
		EventType:      "friend.accepted",
		EventVersion:   1,
		OccurredAt:     time.Now().UTC(),
		Producer:       "social-graph-service",
		IdempotencyKey: idempotencyKey,
		Payload: FriendAcceptedPayload{
			RequestID:  reqID.String(),
			FromUserID: fromUserID.String(),
			ToUserID:   toUserID.String()}}
	payload, err := json.Marshal(envelope)
	if err != nil {
		return domain.OutboxEvent{}, err
	}
	return domain.OutboxEvent{Topic: "socset.friends.events", PartitionKey: fromUserID.String(), Payload: payload}, nil
}
