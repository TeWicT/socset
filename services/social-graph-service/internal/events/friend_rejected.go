package events

import (
	"encoding/json"
	"social-graph-service/internal/domain"
	"time"

	"github.com/google/uuid"
)

type FriendRejectedPayload struct {
	FromUserID string `json:"from_user_id"`
	ToUserID   string `json:"to_user_id"`
}

func BuildFriendRejected(reqID, fromUserID, toUserID uuid.UUID) (domain.OutboxEvent, error) {
	idempotencyKey := "friend.rejected:" + reqID.String()
	envelope := Envelope[FriendRejectedPayload]{
		EventType:      "friend.rejected",
		EventVersion:   1,
		OccurredAt:     time.Now().UTC(),
		Producer:       "social-graph-service",
		IdempotencyKey: idempotencyKey,
		Payload: FriendRejectedPayload{
			FromUserID: fromUserID.String(),
			ToUserID:   toUserID.String()}}
	payload, err := json.Marshal(envelope)
	if err != nil {
		return domain.OutboxEvent{}, err
	}
	return domain.OutboxEvent{Topic: "socset.friends.events", PartitionKey: fromUserID.String(), Payload: payload}, nil
}
