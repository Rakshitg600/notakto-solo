package fxns

import "github.com/rakshitg600/notakto-solo/types"
import "github.com/rakshitg600/notakto-solo/sessions"
import "github.com/google/uuid"

func generateSessionId() string {
	return uuid.New().String()
}


func CreateGame( NumberOfBoards types.BoardNumber, BoardSize types.BoardSize, Difficulty types.DifficultyLevel )types.GameResponse{
	// Initialize boards
	initialBoards := make([]types.BoardState, int (NumberOfBoards))
	for i := 0; i < int(NumberOfBoards); i++ {
		board := make([]string, int(BoardSize)*int(BoardSize))
		for j := range board {
			board[j] = ""
		}
		initialBoards[i] = board
	}

	sessionId := generateSessionId()
	gameOver := false

	gameState := types.GameState{
		Boards:         initialBoards,
		CurrentPlayer:  1,
		Winner:         "",
		BoardSize:      BoardSize,
		NumberOfBoards: NumberOfBoards,
		Difficulty:     Difficulty,
		GameHistory:    [][]types.BoardState{initialBoards},
		SessionId:      sessionId,
		GameOver:       &gameOver,
	}

	sessions.SetGame(sessionId,gameState)

	return types.GameResponse{
		Success:   true,
		SessionId: sessionId,
		GameState: gameState,
	}
}