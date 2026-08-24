package postgres

import (
	"auth-service/internal/domain"
	"auth-service/internal/repository"
	"context"

	"github.com/google/uuid"
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
VALUES ($1, $2, $3);
`

func (r *OutboxRepo) Insert(ctx context.Context, event domain.OutboxEvent) error {
	return insertOutbox(ctx, r.pool, event)
}
func insertOutbox(ctx context.Context, q execer, event domain.OutboxEvent) error {
	_, err := q.Exec(ctx, insertQuery, event.Topic, event.PartitionKey, event.Payload)
	if err != nil {
		return err
	}
	return nil
}

const findNotPublishedQuery = `
SELECT id, topic, partition_key, payload
FROM outbox_events
WHERE published_at IS NULL
ORDER BY created_at
LIMIT $1;
`

func (r *OutboxRepo) FindNotPublished(ctx context.Context, limit int) ([]domain.OutboxEvent, error) {
	rows, err := r.pool.Query(ctx, findNotPublishedQuery, limit)
	if err != nil {
		return []domain.OutboxEvent{}, err
	}
	defer rows.Close()
	result := make([]domain.OutboxEvent, 0, limit)

	for rows.Next() {
		var id uuid.UUID
		var topic string
		var partition_key string
		var payload []byte
		err := rows.Scan(&id, &topic, &partition_key, &payload)
		if err != nil {
			return []domain.OutboxEvent{}, err
		}
		result = append(result, domain.OutboxEvent{ID: id, Topic: topic, PartitionKey: partition_key, Payload: payload})

	}
	err = rows.Err()
	if err != nil {
		return []domain.OutboxEvent{}, err
	}
	return result, nil
}

const setPublishNowQuery = `
UPDATE outbox_events
SET published_at = NOW()
WHERE id = $1
`

func (r *OutboxRepo) SetPublishNow(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, setPublishNowQuery, id)
	if err != nil {
		return err
	}
	return nil
}
