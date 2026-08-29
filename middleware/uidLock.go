package middleware

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rakshitg600/notakto-solo/contextkey"
	"github.com/rakshitg600/notakto-solo/logic"
	"github.com/redis/go-redis/v9"
)

const (
	lockTTL = 12 * time.Second
)

// UIDLockMiddleware returns middleware that serializes requests per UID using
// a distributed lock in Valkey/Redis. Must run after FirebaseAuthMiddleware.
func UIDLockMiddleware(rdb *redis.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()

			uid, ok := contextkey.UIDFromContext(ctx)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "Missing UID")
			}

			lockKey := "lock:uid:" + uid // The key in key val pairs

			unlock, err := logic.AcquireRedisLock(ctx, rdb, lockKey, lockTTL)
			if errors.Is(err, logic.ErrRedisLockAlreadyHeld) {
				return echo.NewHTTPError(http.StatusTooManyRequests, "Could not acquire lock, try again later")
			}
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Lock service unavailable")
			}

			// Ensure unlock runs after the handler, even on panic.
			// Use context.Background() because the request ctx may be canceled.
			defer func() {
				unlockCtx, unlockCancel := context.WithTimeout(context.Background(), 2*time.Second)
				defer unlockCancel()
				if err := unlock(unlockCtx); err != nil {
					c.Logger().Errorf("uid-lock: failed to unlock %s: %v", lockKey, err)
				}
			}()

			return next(c)
		}
	}
}
