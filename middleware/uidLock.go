package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rakshitg600/notakto-solo/contextkey"
	"github.com/redis/go-redis/v9"
)

// Lua script: delete the key only if the value (request identifier) matches.
// example, reqA acquired lock then got timeout then B acquired and A deleted the key ==> this lua solves this bug
var unlockScript = redis.NewScript(`
if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("DEL", KEYS[1])
end
return 0
`)

const (
	lockTTL       = 10 * time.Second
	lockRetryWait = 50 * time.Millisecond
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

			// value in key value pair, this value is unique to each request ==> uid,requestIdentifier map
			nonce := make([]byte, 16)
			if _, err := rand.Read(nonce); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate lock nonce")
			}
			lockVal := hex.EncodeToString(nonce)

			// Retry acquiring the lock until the request context expires.
			for {
				ok, err := rdb.SetNX(ctx, lockKey, lockVal, lockTTL).Result()
				if err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, "Lock service unavailable")
				}
				if ok {
					break
				}

				select {
				case <-ctx.Done():
					return echo.NewHTTPError(http.StatusTooManyRequests, "Could not acquire lock, try again later")
				case <-time.After(lockRetryWait):
				}
			}

			// Ensure unlock runs after the handler, even on panic.
			defer unlockScript.Run(ctx, rdb, []string{lockKey}, lockVal)

			return next(c)
		}
	}
}
