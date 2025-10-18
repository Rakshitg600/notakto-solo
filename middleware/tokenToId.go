package middleware

import (
	"log"
	"net/http"

	"github.com/rakshitg600/notakto-solo/functions"

	"github.com/labstack/echo/v4"
)

func FirebaseAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing Authorization header")
		}

		idToken := authHeader[len("Bearer "):]
		log.Println("Verifying token:", idToken)
		uid, err := functions.VerifyFirebaseToken(idToken)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
		}
		c.Set("uid", uid)
		return next(c)
	}
}
