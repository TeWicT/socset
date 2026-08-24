package kafka

import (
	"context"

	kafkago "github.com/segmentio/kafka-go"
)

type Producer struct {
	Writer *kafkago.Writer
}

func NewProducer(addr string) ProducerRepo {
	return &Producer{Writer: &kafkago.Writer{Addr: kafkago.TCP(addr)}}
}

type ProducerRepo interface {
	Produce(ctx context.Context, topic, key string, value []byte) error
	Close() error
}

func (p *Producer) Produce(ctx context.Context, topic, key string, value []byte) error {

	msgKey := []byte(key)
	err := p.Writer.WriteMessages(ctx, kafkago.Message{Topic: topic, Key: msgKey, Value: value})
	if err != nil {
		return err
	}
	return nil
}
func (p *Producer) Close() error {
	return p.Writer.Close()
}
