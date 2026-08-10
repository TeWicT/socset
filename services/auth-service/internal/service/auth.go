package service

import (
	"auth-service/internal/domain"
	"auth-service/internal/repository"
	"auth-service/internal/token"
	"context"
	"errors"
	"log"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	users     repository.UserRepo
	sessions  repository.SessionRepo
	jwtSecret string
	validate  *validator.Validate
}

var ErrEmptyData = errors.New("empty data")

func NewAuthService(users repository.UserRepo, sessions repository.SessionRepo, jwtSecret string) *AuthService {
	return &AuthService{users: users, validate: validator.New(), sessions: sessions, jwtSecret: jwtSecret}
}

func (service *AuthService) Register(ctx context.Context, user domain.User) (userID uuid.UUID, accessToken domain.AccessToken, raw string, err error) {
	if user.Email == "" || user.Username == "" || user.Password == "" {
		return uuid.Nil, domain.AccessToken{}, "", ErrEmptyData
	}
	validate := service.validate
	err = validate.Struct(user)
	if err != nil {
		return uuid.Nil, domain.AccessToken{}, "", err
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.Nil, domain.AccessToken{}, "", err
	}
	user.PasswordHash = string(passwordHash)

	userID, err = service.users.Create(ctx, user)
	if err != nil {
		return uuid.Nil, domain.AccessToken{}, "", err
	}
	accessToken, raw, err = service.issueTokens(ctx, userID)
	if err != nil {
		return uuid.Nil, domain.AccessToken{}, "", err
	}
	return userID, accessToken, raw, nil
}

var ErrInvalidCredentials = errors.New("invalid credentials")

func (service *AuthService) Login(ctx context.Context, login, password string) (userID uuid.UUID, accessToken domain.AccessToken, rawRefresh string, err error) {
	if login == "" || password == "" {
		return uuid.Nil, domain.AccessToken{}, "", ErrEmptyData
	}
	user, err := service.users.FindByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return uuid.Nil, domain.AccessToken{}, "", ErrInvalidCredentials
		}
		return uuid.Nil, domain.AccessToken{}, "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return uuid.Nil, domain.AccessToken{}, "", ErrInvalidCredentials
	}
	accessToken, raw, err := service.issueTokens(ctx, user.ID)
	if err != nil {
		return uuid.Nil, domain.AccessToken{}, "", err
	}
	return user.ID, accessToken, raw, nil
}

func (service *AuthService) issueTokens(ctx context.Context, userID uuid.UUID) (accessToken domain.AccessToken, raw string, err error) {
	access, expiresIn, err := token.GenerateAccessToken(userID, service.jwtSecret)

	if err != nil {
		log.Print("err create access token", err, userID)
		return domain.AccessToken{}, "", err
	}
	accessToken = domain.AccessToken{Access: access, ExpiresIn: expiresIn}
	raw, hash, expiresAt, err := token.GenerateRefreshToken()

	if err != nil {
		log.Print("err create refresh token", err, userID)
		return domain.AccessToken{}, "", err
	}
	refreshToken := domain.RefreshToken{Raw: raw, Hash: hash, ExpiresAt: expiresAt}
	err = service.sessions.Create(ctx, userID, refreshToken.Hash, refreshToken.ExpiresAt)
	if err != nil {
		return domain.AccessToken{}, "", err
	}
	return accessToken, refreshToken.Raw, nil
}
