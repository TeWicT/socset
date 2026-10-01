package repository

import (
	"context"
	"social-graph-service/internal/domain"

	"github.com/google/uuid"
)

type OutboxRepo interface {
	FindNotPublished(ctx context.Context, count int) ([]domain.OutboxEvent, error)
	SetPublishNow(ctx context.Context, eventID uuid.UUID) error
}
