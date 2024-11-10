package hangman

import (
	"fmt"
	"net/http"
	"reflect"
	"slices"
	"strconv"

	IDGenerator "github.com/geofpwhite/html_games_engine/IDGenerator"
	interfaces "github.com/geofpwhite/html_games_engine/interfaces"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func HangmanRoutes(r *gin.Engine, upgrader *websocket.Upgrader, games map[string]interfaces.Game, playerHashes map[string]*websocket.Conn, inputChannel chan interfaces.Input) {
	r.Static("/hangman_game/", "./build_hangman/")
	r.GET("/hangman/ws/:gameID", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		gameID, b := c.Params.Get("gameID")
		if err != nil || !b {
			panic("/hangman/ws/:gameID gave an error")
		}
		fmt.Println(games[gameID])
		fmt.Println(gameID)
		handleWebSocketHangman(conn, inputChannel, games[gameID], false, "", playerHashes)
	})

	r.GET("/hangman/new_game", func(c *gin.Context) {
		gState := newGameHangman()
		var game interfaces.Game = gState
		games[gState.gameID] = game
		// newTickerInputChannel := make(chan (inputInfo))
		// tickerInputChannels[gState.gameID] = newTickerInputChannel
		// go (*gState).runTicker(tickerTimeoutChannel, newTickerInputChannel, closeGameChannel)
		c.JSON(200, struct {
			GameID string `json:"gameID"`
		}{GameID: gState.gameID})
	})
	r.GET("/hangman/get_games", func(c *gin.Context) {
		c.String(http.StatusOK, "0")
	})

	r.GET("/hangman/reconnect/:playerHash/:gameID", func(c *gin.Context) {

		playerHash, b := c.Params.Get("playerHash")
		if !b {
			return
		}
		gameID, b := c.Params.Get("gameID")
		if !b {
			return
		}
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println(err)
			conn.Close()
			return
		}

		if games[gameID] != nil {
			handleWebSocketHangman(conn, inputChannel, games[gameID], true, playerHash, playerHashes)
		} else {
			conn.WriteJSON(hangmanClientState{Hash: "undefined", Warning: "1"})
		}
	})

	r.GET("/hangman/valid/:playerHash", func(c *gin.Context) {
		hash, b := c.Params.Get("playerHash")
		if !b {
			return
		}
		if playerHashes[hash] == nil {
			c.String(http.StatusOK, "-1")
		} else {
			var gameID string
			for i, g := range games {
				if reflect.TypeOf(g) == reflect.TypeOf(&hangman{}) {
					for _, p := range (g).(*hangman).Players() {
						if p.PlayerID == hash {
							gameID = i
						}
					}
				}
			}
			c.String(http.StatusOK, gameID)
		}
	})

	r.GET("hangman/exit_game/:playerHash/:gameID", func(c *gin.Context) {
		defer c.String(http.StatusOK, "ok")
		playerHash, _ := c.Params.Get("playerHash")
		gameID, _ := c.Params.Get("gameID")
		_player := playerHashes[playerHash]
		if _player == nil || games[gameID] == nil {
			return
		}
		playerIndex := slices.IndexFunc((games[gameID]).(*hangman).players, func(p *interfaces.Player) bool { return p.PlayerID == playerHash })

		delete(playerHashes, playerHash)
		inputChannel <- &exitGameInput{gameID, playerIndex}
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
				playerIndex = slices.IndexFunc(gState.players, func(p *interfaces.Player) bool { return p.PlayerID == hash })
				if playerIndex == -1 {
					conn.WriteJSON(hangmanClientState{Hash: "undefined", Warning: "2"})
					conn.Close()
					return
				}
				playerHashes[hash] = conn
			}

		} else {
			playerIndex = len(gState.players)
			playerHash := IDGenerator.GenerateID(32)
			hash = playerHash
			newPlayer := interfaces.Player{Username: "Player " + strconv.Itoa(playerIndex+1), PlayerID: playerHash}
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
				GameID:         gState.gameID,
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
			GameID:         gState.gameID,
			ChatLogs:       gState.chatLogs,
		}

		for i, player := range gState.players {
			currentState.PlayerIndex = i
			playerHashes[player.PlayerID].WriteJSON(currentState)
		}

		for {
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				return
			}

			GameID := gState.gameID
			PlayerIndex := slices.IndexFunc(gState.players, func(p *interfaces.Player) bool {
				return p.PlayerID == hash
			})
			switch messageType {
			case websocket.TextMessage:
				pString := string(p)
				switch pString[:2] {
				case "g:":
					Guess := pString[2:]
					inp := guessInput{gameID: GameID, playerIndex: PlayerIndex, guess: Guess}
					inputChannel <- &inp
				case "u:":
					Username := pString[2:]
					inp := usernameInput{gameID: GameID, playerIndex: PlayerIndex, username: Username}
					inputChannel <- &inp
				case "w:":
					Word := pString[2:]
					inp := newWordInput{gameID: GameID, playerIndex: PlayerIndex, newWord: Word}
					inputChannel <- &inp
				case "c:":
					Chat := pString[2:]
					inp := chatInput{gameID: GameID, playerIndex: PlayerIndex, message: Chat}
					inputChannel <- &inp
				case "r:":
					inp := randomlyChooseWordInput{gameID: GameID, playerIndex: PlayerIndex}
					inputChannel <- &inp

				default:
					continue
				}
			}
		}

	}
}
