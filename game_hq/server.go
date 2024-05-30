package hq

import (
	"accounts"
	"interfaces"
	"net/http"

	"connect4"
	"hangman"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// nastiest part of the system.
func serve(inputChannel chan interfaces.Input, games map[string]*interfaces.Game, playerHashes map[string]*websocket.Conn) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	r.LoadHTMLGlob("game_hq/templates/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home_page.go.tmpl", gin.H{})
	})
	hangman.HangmanRoutes(r, &upgrader, games, playerHashes, inputChannel)
	connect4.Connect4Routes(r, &upgrader, games, playerHashes, inputChannel)
	accountRoutes(r, accounts.NewAccountsGamesHandler())

	r.Run("0.0.0.0:8080")

}
