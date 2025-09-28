package types

type BoardState []string
type BoardSize int 
type BoardNumber int
type DifficultyLevel int

type GameState struct {
	Boards         []BoardState  `json:"boards"`
	CurrentPlayer  int           `json:"currentPlayer"`
	Winner         string        `json:"winner"`
	BoardSize      BoardSize     `json:"boardSize"`
	NumberOfBoards BoardNumber   `json:"numberOfBoards"`
	Difficulty     DifficultyLevel `json:"difficulty"`
	GameHistory    [][]BoardState  `json:"gameHistory"`
	SessionId      string        `json:"sessionId"`
	GameOver       *bool         `json:"gameOver,omitempty"` // skip this key if val is empty
}

type CreateGameRequest struct {
	NumberOfBoards BoardNumber `json:"numberOfBoards"`
	BoardSize      BoardSize `json:"boardSize"`
	Difficulty     DifficultyLevel `json:"difficulty"`
}

type GameResponse struct {
    Success   bool      `json:"success"`
    SessionId string    `json:"sessionId"`
    GameState GameState `json:"gameState"`
}