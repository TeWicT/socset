package postgres

import (
	"context"
	"errors"
	"profile-service/internal/domain"
	"profile-service/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProfileRepo struct {
	pool *pgxpool.Pool
}

func CreateProfileRepo(pool *pgxpool.Pool) repository.ProfileRepo {
	return &ProfileRepo{pool: pool}
}

const createProfileQuery = `
INSERT INTO profiles (user_id, display_name)
VALUES ($1, $2)
ON CONFLICT (user_id) DO NOTHING;
`

func (r *ProfileRepo) Create(ctx context.Context, profile domain.Profile) error {
	_, err := r.pool.Exec(ctx, createProfileQuery, profile.UserID, profile.DisplayName)
	if err != nil {
		return err
	}
	return nil
}

const getByUserIDQuery = `
SELECT * FROM profiles
WHERE user_id = $1
LIMIT 1;
`

func (r *ProfileRepo) GetByUserID(ctx context.Context, userID uuid.UUID) (domain.Profile, error) {
	var profile domain.Profile
	err := r.pool.QueryRow(ctx, getByUserIDQuery, userID).Scan(&profile.UserID, &profile.DisplayName, &profile.Bio, &profile.AvatarKey, &profile.IsPrivate, &profile.CreatedAt, &profile.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Profile{}, domain.ErrUserNotFound
		}
		return domain.Profile{}, err
	}
	return profile, nil
}

const updateProfileQuery = `
UPDATE profiles
SET
  display_name = COALESCE($1, display_name),
  bio          = COALESCE($2, bio),
  avatar_key   = COALESCE($3, avatar_key),
  is_private   = COALESCE($4, is_private),
  updated_at   = NOW()
WHERE user_id = $5
RETURNING user_id,display_name,bio,avatar_key,is_private,created_at,updated_at;
`

func (r *ProfileRepo) Update(ctx context.Context, userID uuid.UUID, patch domain.Patch) (*domain.Profile, error) {
	var profile domain.Profile
	err := r.pool.QueryRow(ctx, updateProfileQuery, patch.DisplayName, patch.Bio, patch.AvatarKey, patch.IsPrivate, userID).Scan(&profile.UserID, &profile.DisplayName, &profile.Bio, &profile.AvatarKey, &profile.IsPrivate, &profile.CreatedAt, &profile.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &profile, nil
}
