package domain

import (
	"time"

	"github.com/google/uuid"
)

const RequestStatusPending = "pending"
const RequestStatusAccepted = "accepted"
const RequestStatusRejected = "rejected"
const RequestStatusCanceled = "canceled"

type FriendRequest struct {
	ID         uuid.UUID
	FromUserID uuid.UUID
	ToUserID   uuid.UUID
	Status     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
