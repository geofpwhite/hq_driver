package connectthedots

import "interfaces"

type connectTheDotsAddEdgeInput struct {
	team        int
	coords      [2]int
	gameHash    string
	playerIndex int
}

func (ctdaei *connectTheDotsAddEdgeInput) GameHash() string {
	return ctdaei.gameHash
}
func (ctdaei *connectTheDotsAddEdgeInput) PlayerIndex() int {
	return ctdaei.playerIndex
}
func (ctdaei *connectTheDotsAddEdgeInput) ChangeState(gameObj interfaces.Game) {
	(gameObj).(*connectTheDots).addEdge(ctdaei.coords, ctdaei.team)
}
