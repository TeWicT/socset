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
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
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
			carrier := propagation.MapCarrier{}
			for _, h := range m.Headers {
				carrier[h.Key] = string(h.Value)
			}

			parentCtx := otel.GetTextMapPropagator().Extract(ctx, carrier)

			tr := otel.Tracer("profile-service-kafka")
			msgCtx, span := tr.Start(
				parentCtx,
				m.Topic+" process",
				trace.WithSpanKind(trace.SpanKindConsumer),
			)

			var envelope Envelope
			err = json.Unmarshal(m.Value, &envelope)
			if err != nil {
				span.RecordError(err)
				span.End()
				log.Println(err)
				continue
			}

			if envelope.EventType != "user.registered" {
				span.End()
				continue
			}

			id, err := uuid.Parse(envelope.Payload.UserID)
			if err != nil {
				span.RecordError(err)
				span.End()
				log.Println(err)
				continue
			}
			err = c.ProfileRepo.Create(msgCtx, domain.Profile{UserID: id, DisplayName: envelope.Payload.Username})
			if err != nil {
				span.RecordError(err)
				span.End()
				log.Println(err)
				continue
			}

			span.End()
		}
	}

}
func (c *Consumer) Close() {
	c.Reader.Close()
}
