package domain

import "time"

type AccessToken struct {
	Access    string
	ExpiresIn int64
}

type RefreshToken struct {
	Raw       string
	Hash      string
	ExpiresAt time.Time
}
