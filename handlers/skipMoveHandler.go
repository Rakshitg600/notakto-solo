package handlers

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rakshitg600/notakto-solo/functions"
	"github.com/rakshitg600/notakto-solo/types"
)

func (h *Handler) SkipMoveHandler(c echo.Context) error {
	// ✅ Get UID
	uid, ok := c.Get("uid").(string)
	if !ok || uid == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized: missing or invalid uid")
	}

	log.Printf("SkipMoveHandler called for uid: %s", uid)
	// ✅ Try binding the body
	var req types.SkipMoveRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	boards, gameOver, winner, coinsRewarded, xpRewarded, err := functions.EnsureSkipMove(
		c.Request().Context(),
		h.Queries,
		uid,
		req.SessionID,
	)
	if err != nil {
		c.Logger().Errorf("SkipMove failed: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	resp := types.SkipMoveResponse{
		Boards:        boards,
		Gameover:      gameOver,
		Winner:        winner,
		CoinsRewarded: coinsRewarded,
		XpRewarded:    xpRewarded,
	}
	log.Printf("SkipMoveHandler completed for uid: %s, sessionID: %s, gameOver: %v, winner: %v, coinsRewarded: %d, xpRewarded: %d", uid, req.SessionID, gameOver, winner, coinsRewarded, xpRewarded)
	return c.JSON(http.StatusOK, resp)
}
