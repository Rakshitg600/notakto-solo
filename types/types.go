package types

// CreateGameRequest holds creation params
type CreateGameRequest struct {
	NumberOfBoards int32 `json:"numberOfBoards"`
	BoardSize      int32 `json:"boardSize"`
	Difficulty     int32 `json:"difficulty"`
}

type SignInResponse struct {
	Uid        string `json:"uid"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	ProfilePic string `json:"profile_pic"`
}
type FirebaseTokenInfo struct {
	LocalID string `json:"localId"`
	Email   string `json:"email,omitempty"`
	Name    string `json:"displayName,omitempty"`
	Photo   string `json:"photoUrl,omitempty"`
}
