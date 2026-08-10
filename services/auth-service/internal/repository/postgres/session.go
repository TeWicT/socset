package postgres

import (
	"auth-service/internal/repository"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionRepo struct {
	pool *pgxpool.Pool
}

func CreateSessionRepo(pool *pgxpool.Pool) repository.SessionRepo {
	return &SessionRepo{pool: pool}
}

const createRefreshQuery = `
INSERT INTO refresh_sessions (user_id,token_hash,expires_at)
VALUES ($1, $2, $3)
`

func (r *SessionRepo) Create(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	_, err := r.pool.Exec(ctx, createRefreshQuery, userID, tokenHash, expiresAt)
	if err != nil {
		return err
	}
	return nil
}
