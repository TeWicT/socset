package postgres

import (
	"auth-service/internal/domain"
	"auth-service/internal/repository"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func CreateUserRepo(pool *pgxpool.Pool) repository.UserRepo {
	return &UserRepo{pool: pool}
}

const createUserQuery = `
INSERT INTO users (email,username,password_hash)
VALUES ($1, $2, $3)
RETURNING id
`

func (r *UserRepo) Create(ctx context.Context, user domain.User) (userID uuid.UUID, err error) {
	var id uuid.UUID
	err = r.pool.QueryRow(ctx, createUserQuery, user.Email, user.Username, user.PasswordHash).Scan(&id)
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

const findByLoginQuery = `
SELECT id,email,username,password_hash
FROM users
WHERE email = $1 OR username = $1
LIMIT 1
`

func (r *UserRepo) FindByLogin(ctx context.Context, login string) (domain.User, error) {
	var id uuid.UUID
	var email, username, password_hash string
	err := r.pool.QueryRow(ctx, findByLoginQuery, login).Scan(&id, &email, &username, &password_hash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}
	return domain.User{ID: id, Email: email, Username: username, PasswordHash: password_hash}, nil
}
