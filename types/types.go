package types

// CreateGameRequest holds creation params
type CreateGameRequest struct {
	NumberOfBoards int32 `json:"numberOfBoards"`
	BoardSize      int32 `json:"boardSize"`
	Difficulty     int32 `json:"difficulty"`
}
