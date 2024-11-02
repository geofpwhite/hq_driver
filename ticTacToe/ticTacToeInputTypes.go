package tictactoe

import interfaces "github.com/geofpwhite/html_games_engine/interfaces"

type moveInput struct {
	gameHash    string
	playerIndex int
	x, y, team  int
}

func (mi *moveInput) GameHash() string {
	return mi.gameHash
}
func (mi *moveInput) PlayerIndex() int {
	return mi.playerIndex
}
func (mi *moveInput) ChangeState(gameObj interfaces.Game) {
	if gState, ok := gameObj.(*ticTacToe); ok {
		gState.move(mi.x, mi.y, mi.team)
	}
}
