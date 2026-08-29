package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"

	"firebase.google.com/go/v4/auth"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rakshitg600/notakto-solo/contextkey"
	db "github.com/rakshitg600/notakto-solo/db/generated"
	"github.com/rakshitg600/notakto-solo/store"
)

const ensureLoginMaxRetries = 3

func EnsureLogin(ctx context.Context, pool *pgxpool.Pool, authClient *auth.Client) (profilePic string, name string, email string, username string, isNew bool, err error) {
	uid, ok := contextkey.UIDFromContext(ctx)
	if !ok || uid == "" {
		return "", "", "", "", false, errors.New("missing or invalid uid in context")
	}
	// STEP 1: Try existing session
	queries := db.New(pool)
	existing, err := store.GetPlayerById(ctx, queries)
	if err == nil && existing.Uid != "" {
		name = existing.Name
		email = existing.Email
		username = existing.Username
		if existing.ProfilePic.Valid {
			profilePic = existing.ProfilePic.String
		} else {
			profilePic = ""
		}
		return profilePic, name, email, username, false, nil
	}
	if err == nil && existing.Uid == "" {
		return "", "", "", "", false, errors.New("empty player returned from db")
	}
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return "", "", "", "", false, err
	}
	// STEP 2: Fetch profile from Firebase
	name, email, profilePic, err = GetFirebaseUserProfile(ctx, authClient)
	if err != nil {
		return "", "", "", "", true, err
	}
	signUp, err := loadSignUpConfig(ctx, queries)
	if err != nil {
		return "", "", "", "", true, err
	}
	usernameLists, err := loadUsernameWordLists(ctx, queries)
	if err != nil {
		return "", "", "", "", true, err
	}

	for attempt := 1; attempt <= ensureLoginMaxRetries; attempt++ {
		tx, txErr := pool.BeginTx(ctx, pgx.TxOptions{
			IsoLevel:   pgx.Serializable,
			AccessMode: pgx.ReadWrite,
		})
		if txErr != nil {
			return "", "", "", "", true, txErr
		}

		qtx := queries.WithTx(tx)
		username, err = generateAvailableUsername(ctx, qtx, usernameLists)
		if err == nil {
			err = store.CreatePlayer(ctx, qtx, name, email, profilePic, username)
		}
		if err == nil {
			err = store.CreateWallet(ctx, qtx, signUp)
		}
		if err == nil {
			err = tx.Commit(ctx)
		}
		if err == nil {
			return profilePic, name, email, username, true, nil
		}
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil && !errors.Is(rollbackErr, pgx.ErrTxClosed) {
			log.Printf("EnsureLogin: Failed to rollback transaction: %v", rollbackErr)
		}

		var pgErr *pgconn.PgError
		if !errors.As(err, &pgErr) || (pgErr.Code != "23505" && pgErr.Code != "40001") {
			return "", "", "", "", true, err
		}

		existing, existingErr := store.GetPlayerById(ctx, queries)
		if existingErr == nil && existing.Uid != "" {
			name = existing.Name
			email = existing.Email
			username = existing.Username
			if existing.ProfilePic.Valid {
				profilePic = existing.ProfilePic.String
			} else {
				profilePic = ""
			}
			return profilePic, name, email, username, false, nil
		}
	}

	return "", "", "", "", true, fmt.Errorf("could not create player after %d attempts: %w", ensureLoginMaxRetries, err)
}
