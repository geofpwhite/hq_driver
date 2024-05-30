package hq

import (
	"interfaces"

	"github.com/gorilla/websocket"
)

func outputLoop(outputChannel chan string, games map[string]*interfaces.Game, playerHashes map[string]*websocket.Conn) {
	var game interfaces.Game
	var json interfaces.ClientState
	var conn *websocket.Conn
	for gameHash := range outputChannel {

		game = *games[gameHash]
		json = game.JSON()
		for _, p := range game.Players() {
			conn = playerHashes[(*p).PlayerHash]
			conn.WriteJSON(json)
		}
	}
}
