package domain

import "errors"

var ErrSelfRequest = errors.New("self request")
var ErrAlreadyFriends = errors.New("already friends")
var ErrPendingExist = errors.New("pending exist")
var ErrNoRequest = errors.New("no request")
var ErrNoPending = errors.New("no pending")
var ErrForbidden = errors.New("forbidden")
var ErrNoFriends = errors.New("no friends")
var ErrNoHaveFriends = errors.New("no have friends")
