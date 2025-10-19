package functions

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	db "github.com/rakshitg600/notakto-solo/db/generated"
)

// EnsureSession returns the latest existing session for a user, or creates a new one if none exists.
func EnsureSession(ctx context.Context, q *db.Queries, uid string, numberOfBoards int32, boardSize int32, difficulty int32) (map[string]interface{}, error) {
	// STEP 1: Try existing session
	existing, err := q.GetLatestSessionStateByPlayerId(ctx, uid)
	if err == nil && existing.SessionID != "" {
		// Session found, return it
		return map[string]interface{}{
			"session_id":     existing.SessionID,
			"uid":            existing.Uid,
			"boards":         existing.Boards,
			"current_player": existing.CurrentPlayer.Int32,
			"winner":         existing.Winner.String,
			"board_size":     existing.BoardSize.Int32,
			"number_of_boards": existing.NumberOfBoards.Int32,
			"difficulty":       existing.Difficulty.Int32,
			"game_history":     existing.GameHistory,
			"gameover":         existing.Gameover.Bool,
			"created_at":       existing.CreatedAt.Time,
		}, nil
	}

	// STEP 2: Create a new session
	newSessionID := uuid.New().String()

	// a) Insert into session
	err = q.CreateSession(ctx, db.CreateSessionParams{
		SessionID: newSessionID,
		Uid:       uid,
	})
	if err != nil {
		return nil, err
	}

	// b) Insert initial session state
	err = q.CreateInitialSessionState(ctx, db.CreateInitialSessionStateParams{
		SessionID:      newSessionID,
		Boards:         [][]string{},     // empty, adjust if you need initial board
		BoardSize:      sql.NullInt32{Int32: boardSize, Valid: true},
		NumberOfBoards: sql.NullInt32{Int32: numberOfBoards, Valid: true},
		Difficulty:     sql.NullInt32{Int32: difficulty, Valid: true},
		GameHistory:    [][][]string{},   // empty history
	})
	if err != nil {
		return nil, err
	}

	// STEP 3: Return newly created session state
	return map[string]interface{}{
		"session_id":      newSessionID,
		"uid":             uid,
		"boards":          [][]string{},
		"current_player":  1,
		"winner":          "",
		"board_size":      boardSize,
		"number_of_boards": numberOfBoards,
		"difficulty":      difficulty,
		"game_history":    [][][]string{},
		"gameover":        false,
		"created_at":      time.Now(),
	}, nil
}
