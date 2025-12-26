package functions

import (
	"context"

	db "github.com/rakshitg600/notakto-solo/db/generated"
)

func EnsureGetWallet(ctx context.Context, q *db.Queries, uid string) (
	coins int32,
	xp int32,
	err error,
) {
	wallet, err := q.GetWalletByPlayerId(ctx, uid)
	if err != nil {
		return 0, 0, err
	}
	return wallet.Coins.Int32, wallet.Xp.Int32, nil
}
