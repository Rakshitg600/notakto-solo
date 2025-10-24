package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/rakshitg600/notakto-solo/handlers"
	"github.com/rakshitg600/notakto-solo/middleware"
)

func RegisterRoutes(e *echo.Echo, h *handlers.Handler) {
	e.POST("/create-game", h.CreateGameHandler, middleware.FirebaseAuthMiddleware)
	e.POST("/sign-in", h.SignInHandler, middleware.FirebaseAuthMiddleware)
}
