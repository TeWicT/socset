package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Profile struct {
	UserID      uuid.UUID
	DisplayName string
	Bio         *string
	AvatarKey   *string
	IsPrivate   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

var ErrUserNotFound = errors.New("err profile not found")
