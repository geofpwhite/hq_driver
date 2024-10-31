package hq

import (
	"accounts"
	connectthedots "hq/connectTheDots"
	tictactoe "hq/ticTacToe"
	"html/template"
	"interfaces"
	"net/http"

	connect4 "hq/connect4"
	hangman "hq/hangman"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func mod(a, b int) int {
	return a % b
}

// nastiest part of the system.
func serve(inputChannel chan interfaces.Input, games map[string]interfaces.Game, playerHashes map[string]*websocket.Conn) {

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	r.SetFuncMap(template.FuncMap{"mod": mod})

	r.LoadHTMLGlob("game_hq/templates/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home_page.go.tmpl", gin.H{})
	})
	r.GET("/about", func(c *gin.Context) {
		c.HTML(http.StatusOK, "about.go.tmpl", gin.H{})
	})
	r.GET("/contact", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home_page.go.tmpl", gin.H{})
	})

	hangman.HangmanRoutes(r, &upgrader, games, playerHashes, inputChannel)
	connect4.Connect4Routes(r, &upgrader, games, playerHashes, inputChannel)
	connectthedots.ConnectTheDotsRoutes(r, &upgrader, games, playerHashes, inputChannel)
	tictactoe.TicTacToeRoutes(r, &upgrader, games, playerHashes, inputChannel)
	accountRoutes(r, accounts.NewAccountsGamesHandler())

	r.Run("0.0.0.0:8080")

}
