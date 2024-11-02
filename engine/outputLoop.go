package engine

import (
	interfaces "github.com/geofpwhite/html_games_engine/interfaces"

	"github.com/gorilla/websocket"
)

func OutputLoop(outputChannel <-chan string, games map[string]interfaces.Game, playerHashes map[string]*websocket.Conn) {
	var game interfaces.Game
	var json interfaces.ClientState
	var conn *websocket.Conn
	for gameHash := range outputChannel {

		game = games[gameHash]
		json = game.JSON()
		for _, p := range game.Players() {
			conn = playerHashes[p.PlayerHash]
			if err := conn.WriteJSON(json); err != nil {
				continue
			}
		}
	}
}
