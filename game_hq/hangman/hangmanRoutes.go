package hangman

import (
	"fmt"
	interfaces "hq/interfaces"
	myHash "hq/myHash"
	"net/http"
	"reflect"
	"slices"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func HangmanRoutes(r *gin.Engine, upgrader *websocket.Upgrader, games map[string]interfaces.Game, playerHashes map[string]*websocket.Conn, inputChannel chan interfaces.Input) {
	r.Static("/hangman_game/", "./build_hangman/")
	r.GET("/hangman/ws/:gameHash", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		gameHash, b := c.Params.Get("gameHash")
		if err != nil || !b {
			panic("/hangman/ws/:gameHash gave an error")
		}
		fmt.Println(games[gameHash])
		fmt.Println(gameHash)
		handleWebSocketHangman(conn, inputChannel, games[gameHash], false, "", playerHashes)
	})

	r.GET("/hangman/new_game", func(c *gin.Context) {
		gState := newGameHangman()
		var game interfaces.Game = gState
		games[gState.gameHash] = game
		// newTickerInputChannel := make(chan (inputInfo))
		// tickerInputChannels[gState.gameHash] = newTickerInputChannel
		// go (*gState).runTicker(tickerTimeoutChannel, newTickerInputChannel, closeGameChannel)
		c.JSON(200, struct {
			GameHash string `json:"gameHash"`
		}{GameHash: gState.gameHash})
	})
	r.GET("/hangman/get_games", func(c *gin.Context) {
		c.String(http.StatusOK, "0")
	})

	r.GET("/hangman/reconnect/:playerHash/:gameHash", func(c *gin.Context) {

		playerHash, b := c.Params.Get("playerHash")
		if !b {
			return
		}
		gameHash, b := c.Params.Get("gameHash")
		if !b {
			return
		}
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println(err)
			conn.Close()
			return
		}

		if games[gameHash] != nil {
			handleWebSocketHangman(conn, inputChannel, games[gameHash], true, playerHash, playerHashes)
		} else {
			conn.WriteJSON(hangmanClientState{Hash: "undefined", Warning: "1"})
		}
	})

	r.GET("/hangman/valid/:playerHash", func(c *gin.Context) {
		hash, _ := c.Params.Get("playerHash")
		if playerHashes[hash] == nil {
			c.String(http.StatusOK, "-1")
		} else {
			var gameHash string
			for i, g := range games {
				if reflect.TypeOf(g) == reflect.TypeOf(&hangman{}) {
					for _, p := range (g).(*hangman).Players() {
						if p.PlayerHash == hash {
							gameHash = i
						}
					}
				}
			}
			c.String(http.StatusOK, gameHash)
		}
	})

	r.GET("hangman/exit_game/:playerHash/:gameHash", func(c *gin.Context) {
		defer c.String(http.StatusOK, "ok")
		playerHash, _ := c.Params.Get("playerHash")
		gameHash, _ := c.Params.Get("gameHash")
		_player := playerHashes[playerHash]
		if _player == nil || games[gameHash] == nil {
			return
		}
		playerIndex := slices.IndexFunc((games[gameHash]).(*hangman).players, func(p *interfaces.Player) bool { return p.PlayerHash == playerHash })

		delete(playerHashes, playerHash)
		inputChannel <- &exitGameInput{gameHash, playerIndex}
	})
}
func handleWebSocketHangman(
	conn *websocket.Conn,
	inputChannel chan interfaces.Input,
	gameObj interfaces.Game,
	reconnect bool,
	hash string,
	playerHashes map[string]*websocket.Conn,
) {
	if gState, ok := gameObj.(*hangman); ok {

		var playerIndex int
		if reconnect {
			conn2 := playerHashes[hash]
			if conn2 != nil {
				if err := conn2.Close(); err != nil {
					fmt.Println(err)
				}
				playerIndex = slices.IndexFunc(gState.players, func(p *interfaces.Player) bool { return p.PlayerHash == hash })
				if playerIndex == -1 {
					conn.WriteJSON(hangmanClientState{Hash: "undefined", Warning: "2"})
					conn.Close()
					return
				}
				playerHashes[hash] = conn
			}

		} else {
			playerIndex = len(gState.players)
			playerHash := myHash.Hash(32)
			hash = playerHash
			newPlayer := interfaces.Player{Username: "Player " + strconv.Itoa(playerIndex+1), PlayerHash: playerHash}
			gState.newPlayer(newPlayer)

			playerHashes[playerHash] = conn
			usernames := []string{}
			for _, p := range gState.players {
				usernames = append(usernames, p.Username)
			}
			conn.WriteJSON(hangmanClientState{
				Players:        usernames,
				Turn:           gState.turn,
				Host:           gState.curHostIndex,
				RevealedWord:   gState.revealedWord,
				GuessesLeft:    gState.guessesLeft,
				LettersGuessed: gState.guessed,
				NeedNewWord:    gState.needNewWord,
				Warning:        "",
				PlayerIndex:    playerIndex,
				Winner:         gState.winner,
				GameHash:       gState.gameHash,
				ChatLogs:       gState.chatLogs,
				Hash:           playerHash,
			})

		}
		// gState.connections = append(gState.connections, conn)
		defer conn.Close()
		usernames := []string{}
		for _, p := range gState.players {
			usernames = append(usernames, p.Username)
		}

		currentState := hangmanClientState{
			Players:        usernames,
			Turn:           gState.turn,
			Host:           gState.curHostIndex,
			RevealedWord:   gState.revealedWord,
			GuessesLeft:    gState.guessesLeft,
			LettersGuessed: gState.guessed,
			NeedNewWord:    gState.needNewWord,
			Warning:        "",
			PlayerIndex:    playerIndex,
			Winner:         gState.winner,
			GameHash:       gState.gameHash,
			ChatLogs:       gState.chatLogs,
		}

		for i, player := range gState.players {
			currentState.PlayerIndex = i
			playerHashes[player.PlayerHash].WriteJSON(currentState)
		}

		for {
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				return
			}

			GameHash := gState.gameHash
			PlayerIndex := slices.IndexFunc(gState.players, func(p *interfaces.Player) bool {
				return p.PlayerHash == hash
			})
			switch messageType {
			case websocket.TextMessage:
				pString := string(p)
				switch pString[:2] {
				case "g:":
					Guess := pString[2:]
					inp := guessInput{gameHash: GameHash, playerIndex: PlayerIndex, guess: Guess}
					inputChannel <- &inp
				case "u:":
					Username := pString[2:]
					inp := usernameInput{gameHash: GameHash, playerIndex: PlayerIndex, username: Username}
					inputChannel <- &inp
				case "w:":
					Word := pString[2:]
					inp := newWordInput{gameHash: GameHash, playerIndex: PlayerIndex, newWord: Word}
					inputChannel <- &inp
				case "c:":
					Chat := pString[2:]
					inp := chatInput{gameHash: GameHash, playerIndex: PlayerIndex, message: Chat}
					inputChannel <- &inp
				case "r:":
					inp := randomlyChooseWordInput{gameHash: GameHash, playerIndex: PlayerIndex}
					inputChannel <- &inp

				default:
					continue
				}
			}
		}

	}
}
