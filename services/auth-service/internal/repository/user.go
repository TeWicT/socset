package repository

import (
	"auth-service/internal/domain"
	"context"

	"github.com/google/uuid"
)

type UserRepo interface {
	Create(ctx context.Context, user domain.User) (userID uuid.UUID, err error)
}
