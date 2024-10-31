package hq

import (
	"interfaces"

	"github.com/gorilla/websocket"
)

func Main() {

	var games map[string]interfaces.Game = make(map[string]interfaces.Game)
	var playerHashes map[string]*websocket.Conn = make(map[string]*websocket.Conn)
	inputChannel := make(chan interfaces.Input)
	outputChannel := make(chan string)
	go serve(inputChannel, games, playerHashes)
	go outputLoop(outputChannel, games, playerHashes)
	gameLoop(inputChannel, outputChannel, games)
}
