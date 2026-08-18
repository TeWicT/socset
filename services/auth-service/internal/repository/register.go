package repository

import (
	"auth-service/internal/domain"
	"context"

	"github.com/google/uuid"
)

type RegisterRepo interface {
	CreateUserWithOutbox(
		ctx context.Context,
		user domain.User,
		buildEvent func(userID uuid.UUID) (domain.OutboxEvent, error),
	) (uuid.UUID, error)
}
