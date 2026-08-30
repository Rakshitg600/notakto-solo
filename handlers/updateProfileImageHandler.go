package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rakshitg600/notakto-solo/contextkey"
	"github.com/rakshitg600/notakto-solo/usecase"
	"github.com/rakshitg600/notakto-solo/utils"
)

type UpdateProfileImageRequest struct {
	FileID   string `json:"fileId"`
	FilePath string `json:"filePath"`
}

type UpdateProfileImageResponse struct {
	ProfilePic string `json:"profile_pic"`
	FileID     string `json:"fileId"`
	FilePath   string `json:"filePath"`
}

func (h *Handler) UpdateProfileImageHandler(c echo.Context) error {
	uid, ok := contextkey.UIDFromContext(c.Request().Context())
	if !ok || uid == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized: missing or invalid uid")
	}
	log.Printf("UpdateProfileImageHandler called for uid: %s", uid)

	var req UpdateProfileImageRequest
	if err := utils.DecodeStrictJSON(utils.DecodeStrictJSONParams{Context: c, Dest: &req}); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	result, err := usecase.EnsureUpdateProfileImage(
		c.Request().Context(),
		h.Pool,
		h.ImageKitClient,
		req.FileID,
		req.FilePath,
	)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrInvalidProfileImageRequest):
			return echo.NewHTTPError(http.StatusBadRequest, "fileId and filePath are required, and filePath must match your profile-image upload folder")
		case errors.Is(err, usecase.ErrProfileImagePlayerNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "player profile not found; sign in first")
		case errors.Is(err, usecase.ErrProfileImageFileIDConflict):
			return echo.NewHTTPError(http.StatusConflict, "fileId is already associated with another player")
		default:
			c.Logger().Errorf("UpdateProfileImageHandler failed for uid %s: %v", uid, err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to update profile image")
		}
	}

	log.Printf("Updated profile image for uid %s to %s", uid, result.ProfilePic)
	return c.JSON(http.StatusOK, UpdateProfileImageResponse{
		ProfilePic: result.ProfilePic,
		FileID:     result.FileID,
		FilePath:   result.FilePath,
	})
}
