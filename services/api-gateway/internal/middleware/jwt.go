package middleware

import (
	"api-gateway/internal/denylist"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidToken = errors.New("token is invalid")

func ParseAccess(tokenstring, secret string) (userID string, jti string, exp time.Time, err error) {
	token, err := jwt.ParseWithClaims(tokenstring, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Method.Alg())
		}
		return []byte(secret), nil
	})
	if err != nil {
		return "", "", time.Time{}, ErrInvalidToken
	}
	if !token.Valid {
		return "", "", time.Time{}, ErrInvalidToken
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return "", "", time.Time{}, ErrInvalidToken
	}

	if claims.Subject == "" {
		return "", "", time.Time{}, ErrInvalidToken
	}
	if claims.ExpiresAt == nil {
		return "", "", time.Time{}, ErrInvalidToken
	}
	return claims.Subject, claims.ID, claims.ExpiresAt.Time, nil
}

type ctxKey int

var userIDKey ctxKey = 1
var jtiKey ctxKey = 2
var ttlKey ctxKey = 3

func JWT(secret string, deny *denylist.DenyList) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" {

				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			if !strings.HasPrefix(auth, "Bearer ") {

				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			token := strings.TrimPrefix(auth, "Bearer ")
			userID, jti, exp, err := ParseAccess(token, secret)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			if jti == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), userIDKey, userID)
			ctx = context.WithValue(ctx, jtiKey, jti)
			ctx = context.WithValue(ctx, ttlKey, exp)
			ok, err := deny.IsRevoked(ctx, jti)
			if err != nil || ok {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(userIDKey).(string)
	return id, ok
}
func JtiFromContext(ctx context.Context) (string, bool) {
	jti, ok := ctx.Value(jtiKey).(string)
	return jti, ok
}
func ExpFromContext(ctx context.Context) (time.Time, bool) {
	ttl, ok := ctx.Value(ttlKey).(time.Time)
	return ttl, ok
}
