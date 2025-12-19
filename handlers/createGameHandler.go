package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/rakshitg600/notakto-solo/functions"
	"github.com/rakshitg600/notakto-solo/types"
)

func (h *Handler) CreateGameHandler(c echo.Context) error {
	// ✅ Get UID
	uid, ok := c.Get("uid").(string)
	if !ok || uid == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized: missing or invalid uid")
	}

	// ✅ Try binding the body
	var req types.CreateGameRequest
	if err := c.Bind(&req); err != nil {
		req = types.CreateGameRequest{} // reset if malformed JSON
	}

	// ✅ Apply defaults if fields are zero or invalid
	if req.NumberOfBoards < 1 || req.NumberOfBoards > 5 {
		req.NumberOfBoards = 3
	}
	if req.BoardSize < 2 || req.BoardSize > 5 {
		req.BoardSize = 3
	}
	if req.Difficulty < 1 || req.Difficulty > 5 {
		req.Difficulty = 1
	}
	log.Printf("create game handler called for uid: %s", uid)
	// ✅✅ Logic: get typed values from EnsureSession
	sessionID, uidOut, boards, winner, boardSize, numberOfBoards, difficulty, gameover, createdAt, err := functions.EnsureSession(
		c.Request().Context(),
		h.Queries,
		uid,
		req.NumberOfBoards,
		req.BoardSize,
		req.Difficulty,
	)
	// ✅ Handle errors
	if err != nil {
		c.Logger().Errorf("EnsureSession failed: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	createdAtStr := createdAt.UTC().Format(time.RFC3339)

	resp := types.CreateGameResponse{
		SessionId:      sessionID,
		Uid:            uidOut,
		Boards:         boards,
		Winner:         winner,
		BoardSize:      boardSize,
		NumberOfBoards: numberOfBoards,
		Difficulty:     difficulty,
		Gameover:       gameover,
		CreatedAt:      createdAtStr,
	}

	return c.JSON(http.StatusOK, resp)
}
