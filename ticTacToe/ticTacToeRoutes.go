package tictactoe

import (
	"net/http"
	"strconv"

	interfaces "github.com/geofpwhite/html_games_engine/interfaces"
	myHash "github.com/geofpwhite/html_games_engine/myHash"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func TicTacToeRoutes(r *gin.Engine, upgrader *websocket.Upgrader, games map[string]interfaces.Game, playerHashes map[string]*websocket.Conn, inputChannel chan interfaces.Input) {
	r.GET("/tictactoe", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home_screen_tictactoe.go.tmpl", gin.H{})
	})
	r.GET("/tictactoe/:gameHash", func(c *gin.Context) {
		gameHash, b := c.Params.Get("gameHash")
		if !b {
			return
		}
		c.HTML(http.StatusOK, "tictactoe.go.tmpl", gin.H{"Rows": (games[gameHash]).(*ticTacToe).field})
	})
	r.GET("/tictactoe/ws/:gameHash", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		gameHash, b := c.Params.Get("gameHash")
		if err != nil || !b {
			return
		}
		handleWebSocketTicTacToe(conn, inputChannel, games[gameHash], false, "", playerHashes, gameHash)
	})
	r.GET("/tictactoe/reconnect/:playerHash/:gameHash", func(c *gin.Context) {

	})
	r.GET("/tictactoe/new_game", func(c *gin.Context) {
		gState, hash := NewGameTicTacToe()
		var game interfaces.Game = gState
		games[hash] = game
		c.JSON(200, struct {
			GameHash string `json:"gameHash"`
			Team     int    `json:"team"`
		}{GameHash: hash, Team: 1})
	})

}

func handleWebSocketTicTacToe(conn *websocket.Conn,
	inputChannel chan<- interfaces.Input,
	gameObj interfaces.Game,
	reconnect bool,
	hash string,
	playerHashes map[string]*websocket.Conn,
	gameHash string,
) {
	if gState, ok := gameObj.(*ticTacToe); ok {
		var playerIndex int

		if reconnect {

		} else {
			if gState.playersSize > 1 {
				return
			}
			playerIndex = gState.playersSize
			hash = myHash.Hash(10)
			newPlayer := interfaces.Player{Username: "Player " + strconv.Itoa(playerIndex), PlayerHash: hash}
			playerIndex = gState.newPlayer(newPlayer)
			playerHashes[hash] = conn
		}
		defer conn.Close()
		for {
			messageType, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}
			if messageType == websocket.TextMessage {
				x, _ := strconv.Atoi(string(msg[0]))
				y, _ := strconv.Atoi(string(msg[1]))
				mi := moveInput{gameHash: gameHash, x: x, y: y, team: playerIndex + 1, playerIndex: playerIndex}
				inputChannel <- &mi
			}
		}
	}
}
