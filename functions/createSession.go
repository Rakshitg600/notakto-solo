package functions

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	db "github.com/rakshitg600/notakto-solo/db/generated"
)

// EnsureSession returns the latest existing session for a user, or creates a new one if none exists.
// It now returns typed values instead of a generic map so the handler can compose the JSON response.
func EnsureSession(ctx context.Context, q *db.Queries, uid string, numberOfBoards int32, boardSize int32, difficulty int32) (
	sessionID string,
	uidOut string,
	boards [][]string,
	currentPlayer int32,
	winner string,
	boardSizeOut int32,
	numberOfBoardsOut int32,
	difficultyOut int32,
	gameHistory [][][]string,
	gameover bool,
	createdAt time.Time,
	err error,
) {
	// STEP 1: Try existing session
	existing, err := q.GetLatestSessionStateByPlayerId(ctx, uid)
	if err == nil && existing.SessionID != "" {
		sessionID = existing.SessionID
		uidOut = existing.Uid
		boards = existing.Boards

		// convert sql.Null* to plain Go types with sensible defaults
		if existing.CurrentPlayer.Valid {
			currentPlayer = existing.CurrentPlayer.Int32
		} else {
			currentPlayer = 1
		}
		if existing.Winner.Valid {
			winner = existing.Winner.String
		} else {
			winner = ""
		}
		if existing.BoardSize.Valid {
			boardSizeOut = existing.BoardSize.Int32
		} else {
			boardSizeOut = 0
		}
		if existing.NumberOfBoards.Valid {
			numberOfBoardsOut = existing.NumberOfBoards.Int32
		} else {
			numberOfBoardsOut = 0
		}
		if existing.Difficulty.Valid {
			difficultyOut = existing.Difficulty.Int32
		} else {
			difficultyOut = 0
		}
		gameHistory = existing.GameHistory
		if existing.Gameover.Valid {
			gameover = existing.Gameover.Bool
		} else {
			gameover = false
		}
		if existing.CreatedAt.Valid {
			createdAt = existing.CreatedAt.Time
		} else {
			createdAt = time.Time{}
		}

		return sessionID, uidOut, boards, currentPlayer, winner, boardSizeOut, numberOfBoardsOut, difficultyOut, gameHistory, gameover, createdAt, nil
	}

	// STEP 2: Create a new session
	newSessionID := uuid.New().String()

	// a) Insert into session
	if err = q.CreateSession(ctx, db.CreateSessionParams{
		SessionID: newSessionID,
		Uid:       uid,
	}); err != nil {
		return "", "", nil, 0, "", 0, 0, 0, nil, false, time.Time{}, err
	}

	// b) Insert initial session state
	if err = q.CreateInitialSessionState(ctx, db.CreateInitialSessionStateParams{
		SessionID:      newSessionID,
		Boards:         [][]string{}, // empty initial boards
		BoardSize:      sql.NullInt32{Int32: boardSize, Valid: true},
		NumberOfBoards: sql.NullInt32{Int32: numberOfBoards, Valid: true},
		Difficulty:     sql.NullInt32{Int32: difficulty, Valid: true},
		GameHistory:    [][][]string{}, // empty history
	}); err != nil {
		return "", "", nil, 0, "", 0, 0, 0, nil, false, time.Time{}, err
	}

	// STEP 3: Return newly created session state values
	return newSessionID, uid, [][]string{}, 1, "", boardSize, numberOfBoards, difficulty, [][][]string{}, false, time.Now(), nil
}
