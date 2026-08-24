package relay

import (
	"auth-service/internal/kafka"
	"auth-service/internal/repository"
	"context"
	"log"
	"time"
)

type Relay struct {
	outbox   repository.OutboxRepo
	interval time.Duration
	producer kafka.ProducerRepo
}

type RelayRepo interface {
	Run(ctx context.Context)
}

func NewRelay(outbox repository.OutboxRepo, interval time.Duration, producer kafka.ProducerRepo) RelayRepo {
	return &Relay{outbox: outbox, interval: interval, producer: producer}
}

func (r *Relay) Run(ctx context.Context) {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			events, err := r.outbox.FindNotPublished(ctx, 100)
			if err != nil {
				log.Println(err)
				break
			}
			for _, event := range events {
				err = r.producer.Produce(ctx, event.Topic, event.PartitionKey, event.Payload)
				if err != nil {
					log.Println(err)
					break
				}
				log.Printf("Send message to kafka,id:%v,partition_key:%v,Payload:%v,topic:%v", event.ID, event.PartitionKey, string(event.Payload), event.Topic)
				err = r.outbox.SetPublishNow(ctx, event.ID)
				if err != nil {
					log.Println(err)
					break
				}

			}

		}
	}
}
