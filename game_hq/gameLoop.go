package hq

import (
	"interfaces"
	"time"
)

func gameLoop(inputChannel chan interfaces.Input, outputChannel chan string, games map[string]*interfaces.Game) {
	var gameHash string
	var game *interfaces.Game
	lastModified := map[*interfaces.Game]time.Time{}

	cleanupFunction := func() {
		ticker := time.NewTicker(20 * time.Minute)
		defer ticker.Stop()
		var lastTick time.Time = time.Now()
		for interval := range ticker.C {
			for _, game := range games {
				if lastTick.Compare(lastModified[game]) > 0 {
					//close the game
					//hmm
				}
			}
			lastTick = interval
		}
	}
	go cleanupFunction()
	for userInput := range inputChannel {
		gameHash = userInput.GameHash()
		game = games[gameHash]
		userInput.ChangeState(game)
		lastModified[game] = time.Now()
		outputChannel <- gameHash
	}
}
