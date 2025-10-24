package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	db "github.com/rakshitg600/notakto-solo/db/generated"
	"github.com/rakshitg600/notakto-solo/functions"
	"github.com/rakshitg600/notakto-solo/types"
)

type Handler struct {
	Queries *db.Queries
}

func NewHandler(q *db.Queries) *Handler {
	return &Handler{Queries: q}
}

func (h *Handler) SignInHandler(c echo.Context) error {
	// ✅ Get UID
	uid, ok := c.Get("uid").(string)
	if !ok || uid == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized: missing or invalid uid")
	}
	idToken, ok := c.Get("idToken").(string)
	if !ok || uid == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized: missing or invalid token")
	}
	profile_pic, name, email, new, err := functions.EnsureLogin(c.Request().Context(), h.Queries, uid, idToken)
	if err != nil {
		c.Logger().Errorf("EnsurePlayer failed: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	resp := types.SignInResponse{
		Uid:        uid,
		Name:       name,
		Email:      email,
		ProfilePic: profile_pic,
		NewAccount: new,
	}
	return c.JSON(http.StatusOK, resp)
}
