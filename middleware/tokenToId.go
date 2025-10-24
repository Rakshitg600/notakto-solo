package middleware

import (
	"net/http"

	"github.com/rakshitg600/notakto-solo/functions"

	"github.com/labstack/echo/v4"
)

func FirebaseAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
import (
	"log"
	"net/http"
	"strings"
)
		uid, err := functions.VerifyFirebaseToken(idToken)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
		}
		c.Set("uid", uid)
		return next(c)
	}
}
