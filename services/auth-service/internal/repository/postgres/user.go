package postgres

import (
	"auth-service/internal/domain"
	"auth-service/internal/repository"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func CreateUserRepo(pool *pgxpool.Pool) repository.UserRepo {
	return &UserRepo{pool: pool}
}

const CreateQuery = `
INSERT INTO users (email,username,password_hash)
VALUES ($1, $2, $3)
RETURNING id
`

func (r *UserRepo) Create(ctx context.Context, user domain.User) (userID uuid.UUID, err error) {
	var id uuid.UUID
	err = r.pool.QueryRow(ctx, CreateQuery, user.Email, user.Username, user.PasswordHash).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return uuid.Nil, domain.ErrEmailOrUsernameTaken
			}
		}
		return uuid.Nil, err
	}
	return id, nil
}
