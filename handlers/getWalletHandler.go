package handlers

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/rakshitg600/notakto-solo/functions"
	"github.com/rakshitg600/notakto-solo/types"
)

func (h *Handler) GetWalletHandler(c echo.Context) error {
	uid, ok := c.Get("uid").(string)
	if !ok || uid == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized: missing or invalid uid")
	}
	log.Printf("GetWalletHandler called for uid: %s", uid)
	coins, xp, err := functions.EnsureGetWallet(c.Request().Context(), h.Queries, uid)
	if err != nil {
		c.Logger().Errorf("EnsureGetWallet failed: %v", err)
		return c.JSON(http.StatusOK, types.GetWalletResponse{
			Coins:   coins,
			XP:      xp,
			Success: false,
			Error:   err.Error(),
		})
	}

	resp := types.GetWalletResponse{
		Success: true,
		Coins:   coins,
		XP:      xp,
	}
	log.Printf("GetWalletHandler completed for uid: %s", uid)
	return c.JSON(http.StatusOK, resp)
}
