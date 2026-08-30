package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rakshitg600/notakto-solo/contextkey"
	db "github.com/rakshitg600/notakto-solo/db/generated"
	"github.com/rakshitg600/notakto-solo/imagekitservice"
	"github.com/rakshitg600/notakto-solo/store"
)

var ErrProfileImageFileIDConflict = errors.New("profile-image fileId belongs to another user")

type UpdateProfileImageResult struct {
	ProfilePic string
	FileID     string
	FilePath   string
}

func EnsureUpdateProfileImage(
	ctx context.Context,
	pool *pgxpool.Pool,
	imagekitClient *imagekitservice.Client,
	fileID string,
	filePath string,
) (UpdateProfileImageResult, error) {
	uid, ok := contextkey.UIDFromContext(ctx)
	if !ok || uid == "" {
		return UpdateProfileImageResult{}, errors.New("missing or invalid uid in context")
	}
	if strings.TrimSpace(fileID) == "" || strings.TrimSpace(fileID) != fileID ||
		strings.TrimSpace(filePath) == "" || strings.TrimSpace(filePath) != filePath {
		return UpdateProfileImageResult{}, ErrInvalidProfileImageRequest
	}
	if imagekitClient == nil {
		return UpdateProfileImageResult{}, errors.New("ImageKit client is required")
	}
	if err := imagekitClient.ValidateProfileImageFilePath(uid, filePath); err != nil {
		return UpdateProfileImageResult{}, fmt.Errorf("%w: %v", ErrInvalidProfileImageRequest, err)
	}
	profilePic, err := imagekitClient.ProfileImageURL(filePath)
	if err != nil {
		return UpdateProfileImageResult{}, err
	}

	queries := db.New(pool)
	tx, err := pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.Serializable,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return UpdateProfileImageResult{}, err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Printf("EnsureUpdateProfileImage: Failed to rollback transaction: %v", err)
		}
	}(tx, ctx)

	qtx := queries.WithTx(tx)
	if _, err := store.GetPlayerById(ctx, qtx); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return UpdateProfileImageResult{}, ErrProfileImagePlayerNotFound
		}
		return UpdateProfileImageResult{}, fmt.Errorf("look up player profile: %w", err)
	}

	if _, err := store.UpsertImageForUID(ctx, qtx, fileID, filePath); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return UpdateProfileImageResult{}, ErrProfileImageFileIDConflict
		}
		return UpdateProfileImageResult{}, fmt.Errorf("upsert profile image: %w", err)
	}

	if _, err := store.UpdatePlayerProfilePic(ctx, qtx, profilePic); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return UpdateProfileImageResult{}, ErrProfileImagePlayerNotFound
		}
		return UpdateProfileImageResult{}, fmt.Errorf("update player profile image: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return UpdateProfileImageResult{}, err
	}

	return UpdateProfileImageResult{
		ProfilePic: profilePic,
		FileID:     fileID,
		FilePath:   filePath,
	}, nil
}
