package connectthedots

import (
	interfaces "hq/interfaces"
	myHash "hq/myHash"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func ConnectTheDotsRoutes(r *gin.Engine, upgrader *websocket.Upgrader, games map[string]interfaces.Game, playerHashes map[string]*websocket.Conn, inputChannel chan interfaces.Input) {
	r.GET("/connect-the-dots", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home_screen_connectTheDots.go.tmpl", gin.H{})
	})
	r.GET("/connect-the-dots/:gameHash", func(c *gin.Context) {
		gameHash, b := c.Params.Get("gameHash")
		if !b {
			panic("no game hash")
		}
		str := "auto"
		for range 14 {
			str += " auto"
		}
		c.HTML(http.StatusOK, "connectTheDots.go.tmpl", gin.H{"Rows": (games[gameHash]).(*connectTheDots).field, "SizeInt": 8, "GridTemplate": str, "SizeGrid": [7]int{}})
	})
	r.GET("/connect-the-dots-test", func(c *gin.Context) {
		str := "auto"
		for range 14 {
			str += " auto"
		}
		c.HTML(http.StatusOK, "connectTheDots.go.tmpl", gin.H{"Rows": [15][15]int{}, "SizeInt": 8, "GridTemplate": str, "SizeGrid": [7]int{}})
	})

	r.GET("/connect-the-dots/new_game", func(c *gin.Context) {
		c4, hash := NewGameConnectTheDots(8)
		var g interfaces.Game = c4
		games[hash] = g
		c.JSON(200, hash)
	})
	r.GET("/connect-the-dots/reconnect/:gameHash/:playerHash", func(c *gin.Context) {})
	r.GET("/connect-the-dots/ws/:gameHash", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		gameHash, b := c.Params.Get("gameHash")
		if err != nil || !b {
			panic("/hangman/ws/:gameHash gave an error")
		}
		gameObj := games[gameHash]
		game := gameObj.(*connectTheDots)
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

		HandleWebSocketConnectTheDots(conn, inputChannel, gameObj.(*connectTheDots), false, playerHash, playerHashes, gameHash)
	})
}

func HandleWebSocketConnectTheDots(conn *websocket.Conn,
	inputChannel chan interfaces.Input,
	game *connectTheDots,
	reconnect bool,
	hash string,
	playerHashes map[string]*websocket.Conn, gameHash string) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			return
		}
		playerIndex := slices.IndexFunc(game.players, func(p *interfaces.Player) bool { return p.PlayerHash == hash })
		switch messageType {
		case websocket.TextMessage:
			pString := string(p)
			switch pString[:2] {
			case "a:":
				coords := [2]int{}
				numStrings := strings.Split(pString[2:], "-")
				numString1, numString2 := numStrings[0], numStrings[1]
				num, _ := strconv.Atoi(numString1)
				coords[0] = num
				num, _ = strconv.Atoi(numString2)
				coords[1] = num

				ctdaei := &connectTheDotsAddEdgeInput{
					team:        playerIndex + 1,
					playerIndex: playerIndex,
					coords:      coords,
					gameHash:    gameHash,
				}
				inputChannel <- ctdaei

			default:
				continue
			}
		}
	}
}
