package grpcserver

import (
	"auth-service/internal/domain"
	authv1 "auth-service/internal/gen/auth/v1"
	"auth-service/internal/service"
	"context"
	"errors"
	"log"

	"github.com/go-playground/validator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	authv1.UnimplementedAuthServiceServer
	Auth *service.AuthService
}

func (s *Server) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	user := domain.User{Email: req.Email, Password: req.Password, Username: req.Username}
	userID, accessToken, raw, err := s.Auth.Register(ctx, user)
	if err != nil {
		if errors.Is(err, domain.ErrEmailOrUsernameTaken) {
			return nil, status.Error(codes.AlreadyExists, "email or username busy")
		}
		if errors.Is(err, service.ErrEmptyData) {
			return nil, status.Error(codes.InvalidArgument, "empty data")
		}
		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			return nil, status.Error(codes.InvalidArgument, "invalid data")
		}
		log.Printf("register failed: %v", err)
		return nil, status.Error(codes.Internal, "internal error")
	}
	res := &authv1.RegisterResponse{UserId: userID.String(), AccessToken: accessToken.Access, RefreshToken: raw, ExpiresIn: accessToken.ExpiresIn}
	return res, nil
}

func (s *Server) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	userID, access, raw, err := s.Auth.Login(ctx, req.Login, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrEmptyData) {
			return nil, status.Error(codes.InvalidArgument, "empty data")
		}
		if errors.Is(err, service.ErrInvalidCredentials) {
			return nil, status.Error(codes.Unauthenticated, "invalid credentials")
		}
		log.Printf("login failed: %v", err)
		return nil, status.Error(codes.Internal, "internal error")
	}
	res := &authv1.LoginResponse{UserId: userID.String(), AccessToken: access.Access, RefreshToken: raw, ExpiresIn: access.ExpiresIn}
	return res, nil
}

func (s *Server) Refresh(ctx context.Context, req *authv1.RefreshRequest) (*authv1.RefreshResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (s *Server) Logout(ctx context.Context, req *authv1.LogoutRequest) (*authv1.LogoutResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
