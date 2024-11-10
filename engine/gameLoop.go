package engine

import (
	"time"

	interfaces "github.com/geofpwhite/html_games_engine/interfaces"
)

func GameLoop(inputChannel <-chan interfaces.Input, outputChannel chan<- string, games map[string]interfaces.Game) {
	var gameID string
	var game interfaces.Game
	lastModified := map[interfaces.Game]time.Time{}
	cleanupFunction := func() {
		ticker := time.NewTicker(20 * time.Minute)
		defer ticker.Stop()
		var lastTick time.Time = time.Now()
		for interval := range ticker.C {
			for id, game := range games {
				if lastTick.Compare(lastModified[game]) > 0 {
					//close the game
					//hmm
					delete(games, id)

				}
			}
			lastTick = interval
		}
	}
	go cleanupFunction()
	for userInput := range inputChannel {
		gameID = userInput.GameID()
		game = games[gameID]
		go func(userInput interfaces.Input) {
			userInput.ChangeState(game)
			lastModified[game] = time.Now()
			outputChannel <- gameID
		}(userInput)
	}
}
