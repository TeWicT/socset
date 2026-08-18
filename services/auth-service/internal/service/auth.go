package service

import (
	"auth-service/internal/domain"
	"auth-service/internal/events"
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
	users    repository.UserRepo
	sessions repository.SessionRepo

	jwtSecret string
	validate  *validator.Validate
	register  repository.RegisterRepo
}

var ErrEmptyData = errors.New("empty data")

func NewAuthService(users repository.UserRepo, sessions repository.SessionRepo, register repository.RegisterRepo, jwtSecret string) *AuthService {
	return &AuthService{users: users, validate: validator.New(), sessions: sessions, jwtSecret: jwtSecret, register: register}
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

	userID, err = service.register.CreateUserWithOutbox(ctx, user, func(id uuid.UUID) (domain.OutboxEvent, error) {
		return events.BuildUserRegistered(id, user.Username, user.Email)
	})
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

var ErrInvalidRefresh = errors.New("invalid refresh token")

func (service *AuthService) Refresh(ctx context.Context, raw string) (access domain.AccessToken, newRaw string, err error) {
	if raw == "" {
		return domain.AccessToken{}, "", ErrEmptyData
	}
	hash := token.HashRefreshToken(raw)
	session, err := service.sessions.FindByHash(ctx, hash)
	if err != nil {
		if errors.Is(err, domain.ErrSessionNotFound) {
			return domain.AccessToken{}, "", ErrInvalidRefresh
		}
		return domain.AccessToken{}, "", err
	}
	err = service.sessions.Revoke(ctx, session.ID)
	if err != nil {
		return domain.AccessToken{}, "", err
	}
	access, newRaw, err = service.issueTokens(ctx, session.UserID)
	if err != nil {
		return domain.AccessToken{}, "", err
	}
	return access, newRaw, nil
}

func (service *AuthService) Logout(ctx context.Context, raw string) error {
	if raw == "" {
		return ErrEmptyData
	}
	hash := token.HashRefreshToken(raw)
	session, err := service.sessions.FindByHash(ctx, hash)
	if err != nil {
		if errors.Is(err, domain.ErrSessionNotFound) {
			return nil
		}
		return err

	}

	return service.sessions.Revoke(ctx, session.ID)
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
