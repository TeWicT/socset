package repository

import (
	"context"
	"social-graph-service/internal/domain"

	"github.com/google/uuid"
)

type GraphRepo interface {
	CreateRequestWithOutbox(ctx context.Context, req domain.FriendRequest, event domain.OutboxEvent) (*domain.FriendRequest, error)
	AcceptRequestWithOutbox(ctx context.Context, userID, reqID uuid.UUID, makeEventFunc func(reqID, fromUserID, toUserID uuid.UUID) (domain.OutboxEvent, error)) (*domain.FriendRequest, error)
	RejectRequestWithOutbox(ctx context.Context, userID, reqID uuid.UUID, makeEventFunc func(reqID, fromUserID, toUserID uuid.UUID) (domain.OutboxEvent, error)) (*domain.FriendRequest, error)
	CancelRequestWithOutbox(ctx context.Context, userID, reqID uuid.UUID, makeEventFunc func(reqID, fromUserID, toUserID uuid.UUID) (domain.OutboxEvent, error)) (*domain.FriendRequest, error)
	RemoveFriendshipWithOutbox(ctx context.Context, userIDA, userIDB uuid.UUID, event domain.OutboxEvent) error
	ListFriends(ctx context.Context, userID uuid.UUID) ([]domain.Friend, error)
	ListIncoming(ctx context.Context, userID uuid.UUID) ([]domain.FriendRequest, error)
	ListOutgoing(ctx context.Context, userID uuid.UUID) ([]domain.FriendRequest, error)
	GetStatus(ctx context.Context, userID, toUserID uuid.UUID) (*domain.Relation, error)
}
