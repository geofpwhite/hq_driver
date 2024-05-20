package hq

import "github.com/gorilla/websocket"

func Main() {

	var games map[string]*Game = make(map[string]*Game)
	var playerHashes map[string]*websocket.Conn = make(map[string]*websocket.Conn)
	inputChannel := make(chan Input)
	outputChannel := make(chan string)
	go serve(inputChannel, games, playerHashes)
	go outputLoop(outputChannel, games, playerHashes)
	gameLoop(inputChannel, outputChannel, games)
}
