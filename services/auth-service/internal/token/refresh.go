package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"
)

func GenerateRefreshToken() (raw string, hash string, expiresAt time.Time, err error) {
	buf := make([]byte, 32)
	_, err = rand.Read(buf)
	if err != nil {
		return "", "", time.Time{}, err
	}
	raw = hex.EncodeToString(buf)

	hash = HashRefreshToken(raw)
	expiresAt = time.Now().Add(time.Hour * 7 * 24)
	return raw, hash, expiresAt, nil
}

func HashRefreshToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	hash := hex.EncodeToString(sum[:])
	return hash
}
