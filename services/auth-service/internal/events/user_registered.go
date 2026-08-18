package events

import (
	"auth-service/internal/domain"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Envelope struct {
	EventType      string         `json:"event_type"`
	EventVersion   int            `json:"event_version"`
	OccurredAt     time.Time      `json:"occurred_at"`
	Producer       string         `json:"producer"`
	IdempotencyKey string         `json:"idempotency_key"`
	Payload        UserRegistered `json:"payload"`
}
type UserRegistered struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func BuildUserRegistered(userID uuid.UUID, username, email string) (domain.OutboxEvent, error) {
	idempotencyKey := "user.registered:" + userID.String()
	envelope := Envelope{EventType: "user.registered", EventVersion: 1, OccurredAt: time.Now().UTC(), Producer: "auth-service", IdempotencyKey: idempotencyKey, Payload: UserRegistered{UserID: userID.String(), Email: email, Username: username}}
	payload, err := json.Marshal(envelope)
	if err != nil {
		return domain.OutboxEvent{}, err
	}
	return domain.OutboxEvent{Topic: "socset.users.events", PartitionKey: userID.String(), Payload: payload}, nil
}
