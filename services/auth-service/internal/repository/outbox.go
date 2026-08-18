package repository

import (
	"auth-service/internal/domain"
	"context"
)

type OutboxRepo interface {
	Insert(ctx context.Context, event domain.OutboxEvent) error
}
