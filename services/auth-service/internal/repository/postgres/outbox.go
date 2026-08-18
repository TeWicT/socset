package postgres

import (
	"auth-service/internal/domain"
	"auth-service/internal/repository"
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OutboxRepo struct {
	pool *pgxpool.Pool
}

func CreateOutboxRepo(pool *pgxpool.Pool) repository.OutboxRepo {
	return &OutboxRepo{pool: pool}
}

type execer interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

const insertQuery = `
INSERT INTO outbox_events (topic, partition_key, payload)
VALUES ($1, $2, $3)
`

func (r *OutboxRepo) Insert(ctx context.Context, event domain.OutboxEvent) error {
	return insertOutbox(ctx, r.pool, event)
}
func insertOutbox(ctx context.Context, q execer, event domain.OutboxEvent) error {
	_, err := q.Exec(ctx, insertQuery, event.Topic, event.PartitionKey, event.Payload)
	if err != nil {
		return err
	}
	return errors.New("test rollback")
}
