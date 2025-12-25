package handlers

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/rakshitg600/notakto-solo/functions"
	"github.com/rakshitg600/notakto-solo/types"
)

func (h *Handler) QuitGameHandler(c echo.Context) error {
	uid, ok := c.Get("uid").(string)
	if !ok || uid == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized: missing or invalid uid")
	}
	log.Printf("QuitGameHandler called for uid: %s", uid)
	// ✅ Try binding the body
	var req types.QuitGameRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	success, err := functions.EnsureQuitGame(c.Request().Context(), h.Queries, uid, req.SessionID)
	if err != nil {
		c.Logger().Errorf("EnsureQuitGame failed: %v", err)
		return c.JSON(http.StatusOK, types.QuitGameResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	resp := types.QuitGameResponse{
		Success: success,
	}
	log.Printf("QuitGameHandler completed for uid: %s, success: %v", uid, success)
	return c.JSON(http.StatusOK, resp)
}
