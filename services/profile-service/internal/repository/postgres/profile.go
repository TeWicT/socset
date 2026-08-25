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
