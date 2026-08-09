package domain

import (
	"errors"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Password     string `validate:"required,min=8,max=32"`
	Email        string `validate:"required,email"`
	Username     string `validate:"required,min=6,max=32"`
	PasswordHash string
}

var ErrEmailOrUsernameTaken = errors.New("email or username busy")
