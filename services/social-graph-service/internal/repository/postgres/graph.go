package postgres

import (
	"context"
	"social-graph-service/internal/domain"
	"social-graph-service/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GraphRepo struct {
	pool *pgxpool.Pool
}

func CreateGraphRepo(pool *pgxpool.Pool) repository.GraphRepo {
	return &GraphRepo{pool: pool}
}

const canRequestQuery = `
SELECT
  EXISTS (
    SELECT 1 FROM friendships
    WHERE user_id_a = $1 AND user_id_b = $2
 ) AS are_friends,
  (
    SELECT id FROM friend_requests
    WHERE status = 'pending'
      AND (
        (from_user_id = $1 AND to_user_id = $2)
        OR (from_user_id = $2 AND to_user_id = $1)
      )
    LIMIT 1
  ) AS pending_id;
`

const insertRequestQuery = `
INSERT INTO friend_requests (id,from_user_id, to_user_id,status)
VALUES($1,$2,$3,$4)
RETURNING created_at,updated_at
`

func (gr *GraphRepo) CreateRequestWithOutbox(ctx context.Context, req domain.FriendRequest, event domain.OutboxEvent) (domain.FriendRequest, error) {
	tx, err := gr.pool.Begin(ctx)
	if err != nil {
		return domain.FriendRequest{}, err
	}
	defer tx.Rollback(ctx)

	//Проверка можно ли отправлять req
	userA, userB, err := domain.CanonicalPair(req.FromUserID, req.ToUserID)
	if err != nil {
		return domain.FriendRequest{}, err
	}
	var are_friends bool
	var pending_id *uuid.UUID
	err = tx.QueryRow(ctx, canRequestQuery, userA, userB).Scan(&are_friends, &pending_id)
	if err != nil {
		return domain.FriendRequest{}, err
	}
	if are_friends {
		return domain.FriendRequest{}, domain.ErrAlreadyFriends
	}
	if pending_id != nil {
		return domain.FriendRequest{}, domain.ErrPendingExist
	}
	//Отправка в friend_requests
	err = tx.QueryRow(ctx, insertRequestQuery, req.ID, req.FromUserID, req.ToUserID, domain.RequestStatusPending).Scan(&req.CreatedAt, &req.UpdatedAt)
	if err != nil {
		return domain.FriendRequest{}, err
	}
	err = insertOutbox(ctx, tx, event)
	if err != nil {
		return domain.FriendRequest{}, err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return domain.FriendRequest{}, err
	}
	req.Status = domain.RequestStatusPending
	return req, nil

}

const insertQuery = `
INSERT INTO outbox_events (topic, partition_key, payload)
VALUES ($1, $2, $3);
`

type execer interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

func insertOutbox(ctx context.Context, q execer, event domain.OutboxEvent) error {
	_, err := q.Exec(ctx, insertQuery, event.Topic, event.PartitionKey, event.Payload)
	if err != nil {
		return err
	}
	return nil
}
