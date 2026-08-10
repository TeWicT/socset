package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type SessionRepo interface {
	Create(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error
}
