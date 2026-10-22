package kafka

import (
	"context"

	kafkago "github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
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
	tr := otel.Tracer("auth-service-kafka")
	ctx, span := tr.Start(ctx, topic+" publish", trace.WithSpanKind(trace.SpanKindProducer))
	defer span.End()

	carrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, carrier)

	var headers []kafkago.Header
	for k, v := range carrier {
		headers = append(headers, kafkago.Header{
			Key:   k,
			Value: []byte(v),
		})
	}

	msgKey := []byte(key)
	err := p.Writer.WriteMessages(ctx, kafkago.Message{Topic: topic, Key: msgKey, Value: value, Headers: headers})
	if err != nil {
		span.RecordError(err)
		return err
	}
	return nil
}
func (p *Producer) Close() error {
	return p.Writer.Close()
}
