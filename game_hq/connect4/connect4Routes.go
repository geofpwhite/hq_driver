package connect4

import (
	"fmt"
	"myHash"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"interfaces"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func Connect4Routes(r *gin.Engine, upgrader *websocket.Upgrader, games map[string]interfaces.Game, playerHashes map[string]*websocket.Conn, inputChannel chan interfaces.Input) {
	r.GET("/connect4/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home_screen_connect4.go.tmpl", gin.H{})
	})
	r.GET("/connect4/ws/:gameHash", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		gameHash, b := c.Params.Get("gameHash")
		if err != nil || !b {
			panic("/hangman/ws/:gameHash gave an error")
		}
		gameObj := games[gameHash]
		game := gameObj.(*connect4)
		playerHash := myHash.Hash(10)
		playerHashes[playerHash] = conn
		if game.playersConnected >= 2 {
			//don't let them join
			return
		} else if game.playersConnected == 1 {
			game.players = append(game.players, &interfaces.Player{PlayerHash: playerHash, GameHash: gameHash, PlayerIndex: 1})
			game.playersConnected++
		} else if game.playersConnected == 0 {
			game.players = append(game.players, &interfaces.Player{PlayerHash: playerHash, GameHash: gameHash, PlayerIndex: 0})
			game.playersConnected++
		}
		if !b {
			return
		}
		defer func() {
			conn.Close()
			//handle game exit
		}()

		for {

			x, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}
			if x == websocket.TextMessage {
				switch string(msg) {
				case "r":
					// rotate(game)
					c4i := connect4RotateInput{gameHash: gameHash, playerIndex: -1}
					inputChannel <- &c4i
				default:
					// insert(game, team, column)
					msgStrings := strings.Split(string(msg), ",")
					team, _ := strconv.Atoi(msgStrings[0])
					column, _ := strconv.Atoi(msgStrings[1])
					c4i := connect4InsertInput{gameHash: gameHash, team: team, column: column}
					inputChannel <- &c4i
				}
			}
		}

	})
	r.GET("/connect4/new_game", func(c *gin.Context) {
		c4, hash := newGameConnect4()
		var g interfaces.Game = c4
		games[hash] = g
		c.JSON(200, hash)
	})
	r.GET("/connect4/:gameHash", func(c *gin.Context) {
		gameHash, b := c.Params.Get("gameHash")
		colors := map[string]string{
			"1": "blue",
			"2": "red",
		}
		if !b {
			return
		}
		game := (games[gameHash]).(*connect4)
		fmt.Println(game)
		rows := make([][]string, 8)
		for i := range rows {
			rows[i] = make([]string, 8)
		}

		for i := range game.field {
			for j := range game.field[i] {
				rows[i][j] = strconv.Itoa(game.field[i][j])
			}
		}
		slices.Reverse(rows)
		c.HTML(http.StatusOK, "connect4.go.tmpl", gin.H{
			"Rows":   rows,
			"Colors": colors,
		})

	})
}
