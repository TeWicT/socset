package service

import (
	"auth-service/internal/domain"
	"auth-service/internal/repository"
	"context"
	"errors"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	users    repository.UserRepo
	validate *validator.Validate
}

var ErrEmptyData = errors.New("empty data")

func NewAuthService(users repository.UserRepo) *AuthService {
	return &AuthService{users: users, validate: validator.New()}
}

func (service *AuthService) Register(ctx context.Context, user domain.User) (userID uuid.UUID, err error) {
	if user.Email == "" || user.Username == "" || user.Password == "" {
		return uuid.Nil, ErrEmptyData
	}
	validate := service.validate
	err = validate.Struct(user)
	if err != nil {
		return uuid.Nil, err
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.Nil, err
	}
	user.PasswordHash = string(passwordHash)

	userID, err = service.users.Create(ctx, user)
	if err != nil {
		return uuid.Nil, err
	}
	return userID, err
}
