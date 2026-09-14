package denylist

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type DenyList struct {
	RedisClient *redis.Client
}

func CreateDenyList(redisClient *redis.Client) *DenyList {
	return &DenyList{RedisClient: redisClient}
}

func (d *DenyList) Revoke(ctx context.Context, jti string, ttl time.Duration) error {
	err := d.RedisClient.Set(ctx, "access:deny:"+jti, "revoked", ttl).Err()
	if err != nil {
		return err
	}
	return nil
}

func (d *DenyList) IsRevoked(ctx context.Context, jti string) (bool, error) {
	res, err := d.RedisClient.Exists(ctx, "access:deny:"+jti).Result()
	if err != nil {
		return false, err
	}
	if res == 0 {
		return false, nil
	}
	return true, nil
}
