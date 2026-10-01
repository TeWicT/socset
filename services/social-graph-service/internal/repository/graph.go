package repository

import (
	"context"
	"social-graph-service/internal/domain"
)

type GraphRepo interface {
	CreateRequestWithOutbox(ctx context.Context, req domain.FriendRequest, event domain.OutboxEvent) (domain.FriendRequest, error)
}
