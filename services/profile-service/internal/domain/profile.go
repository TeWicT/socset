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

type Patch struct {
	DisplayName *string
	Bio         *string
	AvatarKey   *string
	IsPrivate   *bool
}

var ErrUserNotFound = errors.New("err profile not found")
var ErrNothingToUpdate = errors.New("nothing to update")
