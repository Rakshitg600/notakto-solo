package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rakshitg600/notakto-solo/handlers"
	"github.com/rakshitg600/notakto-solo/types"
)

// TestCreateHandler tests the /create API
func TestCreateHandler(t *testing.T) {

	// -------------------------------
	// 1. Normal POST request test
	// -------------------------------
	t.Run("Valid POST request", func(t *testing.T) {
		reqBody := types.CreateGameRequest{
			NumberOfBoards: types.BoardNumber(3),
			BoardSize:      types.BoardSize(3),
			Difficulty:     types.DifficultyLevel(2),
		}
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/create", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handlers.CreateHandler(rr, req)

		// Check HTTP status
		if rr.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", rr.Code)
		}

		// Decode JSON response
		var resp types.GameResponse
		if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		// Assertions
		if !resp.Success {
			t.Errorf("expected success=true, got false")
		}
		if resp.SessionId == "" {
			t.Errorf("expected non-empty sessionId")
		}
		if types.BoardNumber(len(resp.GameState.Boards)) != reqBody.NumberOfBoards {
			t.Errorf("expected %d boards, got %d", reqBody.NumberOfBoards, len(resp.GameState.Boards))
		}
		if resp.GameState.BoardSize != reqBody.BoardSize {
			t.Errorf("expected boardSize %d, got %d", reqBody.BoardSize, resp.GameState.BoardSize)
		}
		if resp.GameState.Difficulty != reqBody.Difficulty {
			t.Errorf("expected difficulty %d, got %d", reqBody.Difficulty, resp.GameState.Difficulty)
		}
		if resp.GameState.CurrentPlayer != 1 {
			t.Errorf("expected currentPlayer 1, got %d", resp.GameState.CurrentPlayer)
		}
	})

	// -------------------------------
	// 2. Invalid HTTP method test
	// -------------------------------
	t.Run("Invalid methods return 405", func(t *testing.T) {
		// Explicitly declare the slice type without :=
		invalidMethods := []string{http.MethodGet, http.MethodPut, http.MethodDelete, http.MethodPatch}

		// Classic for loop using explicit index
		for index := 0; index < len(invalidMethods); index++ {
			method := invalidMethods[index] // get value manually
			t.Run(method, func(t *testing.T) {
				// Create a fake HTTP request with this method
				req := httptest.NewRequest(method, "/create", nil)

				// Create a recorder to capture the response
				rr := httptest.NewRecorder()

				// Call the handler
				handlers.CreateHandler(rr, req)

				// Check that status is 405 Method Not Allowed
				if rr.Code != http.StatusMethodNotAllowed {
					t.Errorf("case %d (%s): expected 405, got %d", index, method, rr.Code)
				}
			})
		}
	})

	// -------------------------------
	// 3. Invalid JSON test
	// -------------------------------
	t.Run("Invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/create", bytes.NewBuffer([]byte("{bad json")))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		handlers.CreateHandler(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", rr.Code)
		}
	})
}
