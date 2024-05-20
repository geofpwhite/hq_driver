package hq

func gameLoop(inputChannel chan Input, outputChannel chan string, games map[string]*Game) {
	var gameHash string
	var game *Game
	for userInput := range inputChannel {
		gameHash = userInput.GameHash()
		game = games[gameHash]
		userInput.ChangeState(game)
		outputChannel <- gameHash
	}

}
