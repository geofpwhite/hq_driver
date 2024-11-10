package tictactoe

import (
	"fmt"
	"net/http"
	"strconv"

	IDGenerator "github.com/geofpwhite/html_games_engine/IDGenerator"
	interfaces "github.com/geofpwhite/html_games_engine/interfaces"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func TicTacToeRoutes(r *gin.Engine, upgrader *websocket.Upgrader, games map[string]interfaces.Game, playerHashes map[string]*websocket.Conn, inputChannel chan interfaces.Input) {
	r.GET("/tictactoe", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home_screen_tictactoe.go.tmpl", gin.H{})
	})
	r.GET("/tictactoe/:gameID", func(c *gin.Context) {
		gameID, b := c.Params.Get("gameID")
		if !b {
			return
		}
		c.HTML(http.StatusOK, "tictactoe.go.tmpl", gin.H{"Rows": (games[gameID]).(*ticTacToe).field})
	})
	r.GET("/tictactoe/ws/:gameID", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		gameID, b := c.Params.Get("gameID")
		if err != nil || !b {
			return
		}
		handleWebSocketTicTacToe(conn, inputChannel, games[gameID], false, "", playerHashes, gameID)
	})
	r.GET("/tictactoe/reconnect/:playerHash/:gameID", func(c *gin.Context) {

	})
	r.GET("/tictactoe/new_game", func(c *gin.Context) {
		gState, hash := NewGameTicTacToe()
		var game interfaces.Game = gState
		games[hash] = game
		c.JSON(200, struct {
			GameID string `json:"gameID"`
			Team   int    `json:"team"`
		}{GameID: hash, Team: 1})
	})

}

func handleWebSocketTicTacToe(conn *websocket.Conn,
	inputChannel chan<- interfaces.Input,
	gameObj interfaces.Game,
	reconnect bool,
	hash string,
	playerHashes map[string]*websocket.Conn,
	gameID string,
) {
	if gState, ok := gameObj.(*ticTacToe); ok {
		var playerIndex int

		if reconnect {

		} else {
			if gState.playersSize > 1 {
				return
			}
			playerIndex = gState.playersSize
			hash = IDGenerator.GenerateID(10)
			newPlayer := interfaces.Player{Username: "Player " + strconv.Itoa(playerIndex), PlayerID: hash}
			playerIndex = gState.newPlayer(newPlayer)
			playerHashes[hash] = conn

		}
		defer conn.Close()
		ui := &moveInput{gameID: gameID, playerIndex: playerIndex, team: playerIndex + 1}
		for {
			err := conn.ReadJSON(ui)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(ui)
			inputChannel <- ui
		}
	}
}
