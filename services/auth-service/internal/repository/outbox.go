package repository

import (
	"auth-service/internal/domain"
	"context"

	"github.com/google/uuid"
)

type OutboxRepo interface {
	Insert(ctx context.Context, event domain.OutboxEvent) error
	FindNotPublished(ctx context.Context, limit int) ([]domain.OutboxEvent, error)
	SetPublishNow(ctx context.Context, id uuid.UUID) error
}
