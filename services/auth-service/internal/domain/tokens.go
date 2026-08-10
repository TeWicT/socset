package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type AccessToken struct {
	Access    string
	ExpiresIn int64
}

type RefreshToken struct {
	Raw       string
	Hash      string
	ExpiresAt time.Time
}

type Session struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
	RevokedAt time.Time
	UserAgent string
	CreatedAt time.Time
}

var ErrSessionNotFound = errors.New("session not found")
