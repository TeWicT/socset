package service

import (
	"context"
	"profile-service/internal/domain"
	"profile-service/internal/repository"

	"github.com/google/uuid"
)

type ProfileService struct {
	profiles repository.ProfileRepo
}

func CreateProfileService(profiles repository.ProfileRepo) *ProfileService {
	return &ProfileService{profiles: profiles}
}

func (svc *ProfileService) GetProfile(ctx context.Context, userID uuid.UUID) (*domain.Profile, error) {
	profile, err := svc.profiles.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (svc *ProfileService) UpdateProfile(ctx context.Context, userID uuid.UUID, patch domain.Patch) (*domain.Profile, error) {
	if patch.AvatarKey == nil && patch.Bio == nil && patch.DisplayName == nil && patch.IsPrivate == nil {
		return nil, domain.ErrNothingToUpdate
	}
	profile, err := svc.profiles.Update(ctx, userID, patch)
	if err != nil {
		return nil, err
	}
	return profile, nil
}
