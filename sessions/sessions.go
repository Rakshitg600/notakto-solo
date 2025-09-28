package sessions

import "github.com/rakshitg600/notakto-solo/types"

var GameMap = make(map[string]types.GameState)

func SetGame(id string, state types.GameState) {
	GameMap[id] = state
}

func GetGame(id string) (types.GameState, bool) {
	state, ok := GameMap[id]
	return state, ok
}