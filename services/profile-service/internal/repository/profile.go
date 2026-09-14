package repository

import (
	"context"
	"profile-service/internal/domain"

	"github.com/google/uuid"
)

type ProfileRepo interface {
	Create(ctx context.Context, profile domain.Profile) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (domain.Profile, error)
	Update(ctx context.Context, userID uuid.UUID, patch domain.Patch) (*domain.Profile, error)
}
