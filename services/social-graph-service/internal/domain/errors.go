package domain

import "errors"

var ErrSelfRequest = errors.New("self request")
var ErrAlreadyFriends = errors.New("already friends")
var ErrPendingExist = errors.New("pending exist")
