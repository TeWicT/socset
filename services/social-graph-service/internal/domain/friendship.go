package domain

import (
	"bytes"
	"time"

	"github.com/google/uuid"
)

const RelationStatusNone = "none"
const RelationStatusOutgoingRequest = "outgoing_request"
const RelationStatusIncomingRequest = "incoming_request"
const RelationStatusFriends = "friends"

type Relation struct {
	Status string
	ID     *uuid.UUID
}

type Friendship struct {
	UserIDA   uuid.UUID
	UserIDB   uuid.UUID
	CreatedAt time.Time
}

type Friend struct {
	UserID    uuid.UUID
	CreatedAt time.Time
}

func CanonicalPair(a, b uuid.UUID) (userA, userB uuid.UUID, err error) {
	if a == b {
		return uuid.Nil, uuid.Nil, ErrSelfRequest
	}
	if bytes.Compare(a[:], b[:]) < 0 {
		return a, b, nil
	}
	return b, a, nil

}
