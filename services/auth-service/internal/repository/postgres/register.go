package postgres

import (
	"auth-service/internal/domain"
	"auth-service/internal/repository"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RegisterRepo struct {
	pool *pgxpool.Pool
}

func CreateRegisterRepo(pool *pgxpool.Pool) repository.RegisterRepo {
	return &RegisterRepo{pool: pool}
}

func (r *RegisterRepo) CreateUserWithOutbox(ctx context.Context, user domain.User, buildEvent func(userID uuid.UUID) (domain.OutboxEvent, error),
) (uuid.UUID, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	defer tx.Rollback(ctx)
	userID, err := insertUser(ctx, tx, user)
	if err != nil {
		return uuid.Nil, err
	}
	event, err := buildEvent(userID)
	if err != nil {
		return uuid.Nil, err
	}
	if err := insertOutbox(ctx, tx, event); err != nil {
		return uuid.Nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, err
	}
	return userID, nil
}
