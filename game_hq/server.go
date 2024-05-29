package hq

import (
	"accounts"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// nastiest part of the system.
func serve(inputChannel chan Input, games map[string]*Game, playerHashes map[string]*websocket.Conn) {
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
	hangmanRoutes(r, &upgrader, games, playerHashes, inputChannel)
	connect4Routes(r, &upgrader, games, playerHashes, inputChannel)
	accountRoutes(r, accounts.NewAccountsGamesHandler())

	r.Run("0.0.0.0:8080")

}
