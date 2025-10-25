package functions

import (
	"context"

	db "github.com/rakshitg600/notakto-solo/db/generated"
)

func EnsureUpdateName(ctx context.Context, q *db.Queries, name string, uid string) (string, error) {
	player, err := q.UpdatePlayerName(ctx, db.UpdatePlayerNameParams{
		Uid:  uid,
		Name: name,
	})
	if err != nil {
		return "", err
	}
	return player.Name, nil
}
