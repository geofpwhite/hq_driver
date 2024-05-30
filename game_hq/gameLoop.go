package hq

import "interfaces"

func gameLoop(inputChannel chan interfaces.Input, outputChannel chan string, games map[string]*interfaces.Game) {
	var gameHash string
	var game *interfaces.Game
	for userInput := range inputChannel {
		gameHash = userInput.GameHash()
		game = games[gameHash]
		userInput.ChangeState(game)
		outputChannel <- gameHash
	}

}
