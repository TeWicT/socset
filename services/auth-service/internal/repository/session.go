package repository

import (
	"auth-service/internal/domain"
	"context"
	"time"

	"github.com/google/uuid"
)

type SessionRepo interface {
	Create(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error
	FindByHash(ctx context.Context, tokenHash string) (session domain.Session, err error)
	Revoke(ctx context.Context, sessionID uuid.UUID) error
}
