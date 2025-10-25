package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/rakshitg600/notakto-solo/functions"
	"github.com/rakshitg600/notakto-solo/types"
)

func (h *Handler) UpdateNameHandler(c echo.Context) error {
	// ✅ Get UID
	uid, ok := c.Get("uid").(string)
	if !ok || uid == "" {
		return echo.NewHTTPError(401, "unauthorized: missing or invalid uid")
	}

	// ✅ Try binding the body
	var req types.UpdatePlayerNameRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(400, "invalid request body")
	}
	if req.Name == "" {
		return echo.NewHTTPError(400, "name is required")
	}
	// ✅ Update the name
	updatedName, err := functions.EnsureUpdateName(c.Request().Context(), h.Queries, req.Name, uid)
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}
	// ✅ Return the updated name
	return c.JSON(200, map[string]string{"name": updatedName})
}
