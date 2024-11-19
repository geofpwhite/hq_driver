package main

import (
	engine "github.com/geofpwhite/html_games_engine/engine"
	interfaces "github.com/geofpwhite/html_games_engine/interfaces"

	"github.com/gorilla/websocket"
)

func main() {

	var games map[string]interfaces.Game = make(map[string]interfaces.Game)
	var playerHashes map[string]*websocket.Conn = make(map[string]*websocket.Conn)
	inputChannel := make(chan interfaces.Input)
	outputChannel := make(chan string)
	go engine.Serve(inputChannel, games, playerHashes)
	go engine.OutputLoop(outputChannel, games, playerHashes)

	engine.GameLoop(inputChannel, outputChannel, games)
}
