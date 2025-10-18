package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	db "github.com/rakshitg600/notakto-solo/db/generated"
)

type Handler struct {
	Queries *db.Queries
}

func NewHandler(q *db.Queries) *Handler {
	return &Handler{Queries: q}
}
func (h *Handler) CreateHandler(c echo.Context) error {
	uidVal := c.Get("uid")
	if uidVal == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "missing uid")
	}

	player, err := h.Queries.GetPlayerById(c.Request().Context(), uidVal.(string))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "player not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get player")
	}
	return c.JSON(http.StatusOK, player)
}
