package postgres

import (
	"context"
	"errors"
	"social-graph-service/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type execer interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}
type queryRower interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

const findRequestWithIDQuery = `
SELECT id, from_user_id, to_user_id, status, created_at, updated_at FROM friend_requests
WHERE id = $1;
`

func findRequestByID(ctx context.Context, q queryRower, reqID uuid.UUID) (*domain.FriendRequest, error) {
	var req domain.FriendRequest
	err := q.QueryRow(ctx, findRequestWithIDQuery, reqID).Scan(&req.ID, &req.FromUserID, &req.ToUserID, &req.Status, &req.CreatedAt, &req.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNoRequest
		}
		return nil, err
	}
	return &req, nil
}

const updateFriendRequestStatusQuery = `
UPDATE friend_requests
SET status = $1, updated_at = NOW()
WHERE id = $2 AND status = 'pending'
RETURNING id, from_user_id, to_user_id, status, created_at, updated_at;
`

func updateFriendRequest(ctx context.Context, q queryRower, status string, reqID uuid.UUID) (*domain.FriendRequest, error) {
	var req domain.FriendRequest
	err := q.QueryRow(ctx, updateFriendRequestStatusQuery, status, reqID).Scan(&req.ID, &req.FromUserID, &req.ToUserID, &req.Status, &req.CreatedAt, &req.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNoPending
		}
		return nil, err
	}
	return &req, nil
}

const insertQuery = `
INSERT INTO outbox_events (topic, partition_key, payload)
VALUES ($1, $2, $3);
`

func insertOutbox(ctx context.Context, q execer, event domain.OutboxEvent) error {
	_, err := q.Exec(ctx, insertQuery, event.Topic, event.PartitionKey, event.Payload)
	if err != nil {
		return err
	}
	return nil
}
