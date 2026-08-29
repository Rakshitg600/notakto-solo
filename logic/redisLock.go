package logic

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/rakshitg600/notakto-solo/lua"
	"github.com/redis/go-redis/v9"
)

var ErrRedisLockAlreadyHeld = errors.New("redis lock already held")

func AcquireRedisLock(ctx context.Context, rdb *redis.Client, key string, ttl time.Duration) (func(context.Context) error, error) {
	nonce := make([]byte, 16)
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	lockVal := hex.EncodeToString(nonce)

	ok, err := rdb.SetNX(ctx, key, lockVal, ttl).Result()
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrRedisLockAlreadyHeld
	}

	return func(unlockCtx context.Context) error {
		return lua.Unlock.Run(unlockCtx, rdb, []string{key}, lockVal).Err()
	}, nil
}
