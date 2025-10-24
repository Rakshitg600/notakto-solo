package functions

import (
	"context"
	"database/sql"

	db "github.com/rakshitg600/notakto-solo/db/generated"
)

// EnsureSession returns the latest existing session for a user, or creates a new one if none exists.
// It returns typed values for the handler to compose the JSON response.
func EnsureLogin(ctx context.Context, q *db.Queries, uid string, idToken string) (profile_pic string, name string, email string, new bool, err error) {
	// STEP 1: Try existing session
	existing, err := q.GetPlayerById(ctx, uid)
	if err == nil && existing.Uid != "" {
		name = existing.Name
		email = existing.Email
		if existing.ProfilePic.Valid {
			profile_pic = existing.ProfilePic.String
		} else {
			profile_pic = ""
		}
		return profile_pic, name, email, false, nil
	}
	// STEP 2: Fetch from Firebase
	uid, name, email, profile_pic, err = VerifyFirebaseToken(idToken)
	if err != nil {
		return "", "", "", true, err
	}

	// STEP 3: Create new player
	err = q.CreatePlayer(ctx, db.CreatePlayerParams{
		Uid:   uid,
		Name:  name,
		Email: email,
		ProfilePic: sql.NullString{
			String: profile_pic,
			Valid:  profile_pic != "",
		},
	})
	if err != nil {
		return "", "", "", true, err
	}
	return profile_pic, name, email, true, nil
}
