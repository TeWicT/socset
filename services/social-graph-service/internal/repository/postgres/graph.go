package postgres

import (
	"context"
	"social-graph-service/internal/domain"
	"social-graph-service/internal/repository"

	"github.com/google/uuid"
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

func (gr *GraphRepo) RejectRequestWithOutbox(ctx context.Context, userID, reqID uuid.UUID, makeEventFunc func(reqID, fromUserID, toUserID uuid.UUID) (domain.OutboxEvent, error)) (*domain.FriendRequest, error) {
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
	if req.Status == domain.RequestStatusRejected {
		return req, nil
	}
	if req.Status != domain.RequestStatusPending {
		return nil, domain.ErrNoPending
	}
	updreq, err := updateFriendRequest(ctx, tx, domain.RequestStatusRejected, reqID)
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

func (gr *GraphRepo) CancelRequestWithOutbox(ctx context.Context, userID, reqID uuid.UUID, makeEventFunc func(reqID, fromUserID, toUserID uuid.UUID) (domain.OutboxEvent, error)) (*domain.FriendRequest, error) {
	tx, err := gr.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	req, err := findRequestByID(ctx, tx, reqID)
	if err != nil {
		return nil, err
	}
	if req.FromUserID != userID {
		return nil, domain.ErrForbidden
	}
	if req.Status == domain.RequestStatusCanceled {
		return req, nil
	}
	if req.Status != domain.RequestStatusPending {
		return nil, domain.ErrNoPending
	}
	updreq, err := updateFriendRequest(ctx, tx, domain.RequestStatusCanceled, reqID)
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
