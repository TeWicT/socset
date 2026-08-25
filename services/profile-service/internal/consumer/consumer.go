package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"profile-service/internal/domain"
	"profile-service/internal/repository"

	"github.com/google/uuid"
	kafkago "github.com/segmentio/kafka-go"
)

type Consumer struct {
	Reader      *kafkago.Reader
	ProfileRepo repository.ProfileRepo
	GroupID     string
}

type ConsumerRepo interface {
	Run(ctx context.Context)
	Close()
}

func CreateConsumer(profileRepo repository.ProfileRepo, groupID, brokersAddr string) ConsumerRepo {
	return &Consumer{Reader: kafkago.NewReader(kafkago.ReaderConfig{Brokers: []string{brokersAddr}, Topic: "socset.users.events", GroupID: groupID}), ProfileRepo: profileRepo, GroupID: groupID}
}

type Envelope struct {
	EventType      string         `json:"event_type"`
	IdempotencyKey string         `json:"idempotency_key"`
	Payload        UserRegistered `json:"payload"`
}

type UserRegistered struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
}

func (c *Consumer) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			m, err := c.Reader.ReadMessage(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					return
				}
				log.Println(err)
				continue
			}
			var envelope Envelope
			err = json.Unmarshal(m.Value, &envelope)
			if err != nil {
				log.Println(err)
				continue
			}
			if envelope.EventType != "user.registered" {
				continue
			}
			id, err := uuid.Parse(envelope.Payload.UserID)
			if err != nil {
				log.Println(err)
				continue
			}
			err = c.ProfileRepo.Create(ctx, domain.Profile{UserID: id, DisplayName: envelope.Payload.Username})
			if err != nil {
				log.Println(err)
				continue
			}
		}
	}

}
func (c *Consumer) Close() {
	c.Reader.Close()
}
