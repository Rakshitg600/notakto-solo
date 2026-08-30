package store

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/rakshitg600/notakto-solo/contextkey"
	db "github.com/rakshitg600/notakto-solo/db/generated"
)

func UpsertImageForUID(ctx context.Context, q *db.Queries, fileID, filePath string) (image db.Image, err error) {
	uid, ok := contextkey.UIDFromContext(ctx)
	if !ok || uid == "" {
		return db.Image{}, errors.New("missing or invalid uid in context")
	}
	start := time.Now()
	image, err = q.UpsertImageForUID(ctx, db.UpsertImageForUIDParams{
		FileID:   fileID,
		Uid:      uid,
		FilePath: filePath,
	})
	if time.Since(start) > 2*time.Second {
		log.Printf("Upsert image for uid took %v, err: %v", time.Since(start), err)
	}
	return image, err
}
