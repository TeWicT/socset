package postgres

import (
	"context"
	"errors"
	"social-graph-service/internal/domain"
	"social-graph-service/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

func (gr *GraphRepo) CreateRequestWithOutbox(ctx context.Context, req domain.FriendRequest, event domain.OutboxEvent) (*domain.FriendRequest, error) {
	tx, err := gr.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	//Проверка можно ли отправлять req
	userA, userB, err := domain.CanonicalPair(req.FromUserID, req.ToUserID)
	if err != nil {
		return nil, err
	}
	var are_friends bool
	var pending_id *uuid.UUID
	err = tx.QueryRow(ctx, canRequestQuery, userA, userB).Scan(&are_friends, &pending_id)
	if err != nil {
		return nil, err
	}
	if are_friends {
		return nil, domain.ErrAlreadyFriends
	}
	if pending_id != nil {
		return nil, domain.ErrPendingExist
	}
	//Отправка в friend_requests
	err = tx.QueryRow(ctx, insertRequestQuery, req.ID, req.FromUserID, req.ToUserID, domain.RequestStatusPending).Scan(&req.CreatedAt, &req.UpdatedAt)
	if err != nil {
		return nil, err
	}
	err = insertOutbox(ctx, tx, event)
	if err != nil {
		return nil, err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}
	req.Status = domain.RequestStatusPending
	return &req, nil

}

const acceptRequestQuery = `
INSERT INTO friendships (user_id_a, user_id_b)
VALUES ($1,$2)
RETURNING user_id_a, user_id_b, created_at;
`

func (gr *GraphRepo) AcceptRequestWithOutbox(ctx context.Context, userID, reqID uuid.UUID, makeEventFunc func(reqID, fromUserID, toUserID uuid.UUID) (domain.OutboxEvent, error)) (*domain.FriendRequest, error) {
	tx, err := gr.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	req, err := findRequestByID(ctx, tx, reqID)
	if err != nil {
		return nil, err
	}
	if req.ToUserID != userID {
		return nil, domain.ErrForbidden
	}
	if req.Status == domain.RequestStatusAccepted {
		return req, nil
	}
	if req.Status != domain.RequestStatusPending {
		return nil, domain.ErrNoPending
	}
	updreq, err := updateFriendRequest(ctx, tx, domain.RequestStatusAccepted, reqID)
	if err != nil {
		return nil, err
	}
	userA, userB, err := domain.CanonicalPair(req.FromUserID, req.ToUserID)
	if err != nil {
		return nil, err
	}
	var friendship domain.Friendship
	err = tx.QueryRow(ctx, acceptRequestQuery, userA, userB).Scan(&friendship.UserIDA, &friendship.UserIDB, &friendship.CreatedAt)
	if err != nil {
		return nil, err
	}

	event, err := makeEventFunc(updreq.ID, updreq.FromUserID, updreq.ToUserID)
	if err != nil {
		return nil, err
	}
	err = insertOutbox(ctx, tx, event)
	if err != nil {
		return nil, err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}
	return updreq, nil
}

const insertQuery = `
INSERT INTO outbox_events (topic, partition_key, payload)
VALUES ($1, $2, $3);
`

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

func insertOutbox(ctx context.Context, q execer, event domain.OutboxEvent) error {
	_, err := q.Exec(ctx, insertQuery, event.Topic, event.PartitionKey, event.Payload)
	if err != nil {
		return err
	}
	return nil
}
