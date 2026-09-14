package grpcserver

import (
	"context"
	"errors"
	"profile-service/internal/domain"
	profilev1 "profile-service/internal/gen/profile/v1"
	"profile-service/internal/service"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	profilev1.UnimplementedProfileServiceServer
	Profiles *service.ProfileService
}

func (s *Server) GetProfile(ctx context.Context, req *profilev1.GetProfileRequest) (*profilev1.ProfileResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	profile, err := s.Profiles.GetProfile(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &profilev1.ProfileResponse{UserId: profile.UserID.String(), DisplayName: profile.DisplayName, Bio: pointerToString(profile.Bio), AvatarKey: pointerToString(profile.AvatarKey), IsPrivate: profile.IsPrivate}, nil
}

func (s *Server) UpdateProfile(ctx context.Context, req *profilev1.UpdateProfileRequest) (*profilev1.ProfileResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	patch := domain.Patch{DisplayName: req.DisplayName, Bio: req.Bio, AvatarKey: req.AvatarKey, IsPrivate: req.IsPrivate}
	profile, err := s.Profiles.UpdateProfile(ctx, userID, patch)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		if errors.Is(err, domain.ErrNothingToUpdate) {
			return nil, status.Error(codes.InvalidArgument, "nothing to update")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &profilev1.ProfileResponse{UserId: profile.UserID.String(), DisplayName: profile.DisplayName, Bio: pointerToString(profile.Bio), AvatarKey: pointerToString(profile.AvatarKey), IsPrivate: profile.IsPrivate}, nil
}

func pointerToString(pointer *string) string {
	if pointer == nil {
		return ""
	}
	return *pointer
}
