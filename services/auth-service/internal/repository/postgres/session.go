package postgres

import (
	"auth-service/internal/domain"
	"auth-service/internal/repository"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

const findByHashQuery = `
SELECT id, user_id, expires_at
FROM refresh_sessions
WHERE token_hash = $1 AND revoked_at IS NULL AND expires_at > now()
LIMIT 1
`

func (r *SessionRepo) FindByHash(ctx context.Context, tokenHash string) (session domain.Session, err error) {

	err = r.pool.QueryRow(ctx, findByHashQuery, tokenHash).Scan(&session.ID, &session.UserID, &session.ExpiresAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Session{}, domain.ErrSessionNotFound
		}
		return domain.Session{}, err
	}
	return session, nil
}

const revokeQuery = `
UPDATE refresh_sessions
SET revoked_at = now()
WHERE id = $1
`

func (r *SessionRepo) Revoke(ctx context.Context, sessionID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, revokeQuery, sessionID)
	if err != nil {
		return err
	}
	return nil
}
