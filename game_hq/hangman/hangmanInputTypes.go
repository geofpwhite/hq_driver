package hangman

import "interfaces"

type usernameInput struct {
	username    string
	gameHash    string
	playerIndex int
}
type newWordInput struct {
	newWord     string
	gameHash    string
	playerIndex int
}
type randomlyChooseWordInput struct {
	gameHash    string
	playerIndex int
}
type guessInput struct {
	guess       string
	gameHash    string
	playerIndex int
}
type chatInput struct {
	message     string
	gameHash    string
	playerIndex int
}
type exitGameInput struct {
	gameHash    string
	playerIndex int
}
type closeGameInput struct {
	gameHash    string
	playerIndex int
}

func (ui *usernameInput) GameHash() string {
	return ui.gameHash
}
func (ui *usernameInput) PlayerIndex() int {
	return ui.playerIndex
}
func (ui *usernameInput) ChangeState(gameObj interfaces.Game) {
	gState := (gameObj).(*hangman)
	gState.changeUsername(ui.playerIndex, ui.username)
}

func (nwi *newWordInput) GameHash() string {
	return nwi.gameHash
}
func (nwi *newWordInput) PlayerIndex() int {
	return nwi.playerIndex
}
func (nwi *newWordInput) ChangeState(gameObj interfaces.Game) {
	gState := (gameObj).(*hangman)
	if gState.needNewWord && nwi.playerIndex == gState.curHostIndex {
		gState.newWord(nwi.newWord)
	}

}
func (gi *guessInput) GameHash() string {
	return gi.gameHash
}
func (gi *guessInput) PlayerIndex() int {
	return gi.playerIndex
}

func (gi *guessInput) ChangeState(gameObj interfaces.Game) {
	gState := (gameObj).(*hangman)
	if gi.playerIndex == gState.turn {
		// fmt.Println("guess")
		gState.guess(rune(gi.guess[0]))
	}
}

func (ci *chatInput) GameHash() string {
	return ci.gameHash
}
func (ci *chatInput) PlayerIndex() int {
	return ci.playerIndex
}
func (ci *chatInput) ChangeState(gameObj interfaces.Game) {
	gState := (gameObj).(*hangman)
	gState.chat(ci.message, ci.playerIndex)
}

func (rcwi *randomlyChooseWordInput) GameHash() string {
	return rcwi.gameHash
}
func (rcwi *randomlyChooseWordInput) PlayerIndex() int {
	return rcwi.playerIndex
}
func (rcwi *randomlyChooseWordInput) ChangeState(gameObj interfaces.Game) {
	gState := (gameObj).(*hangman)
	gState.randomNewWord()
}

func (egi *exitGameInput) GameHash() string {
	return egi.gameHash
}
func (egi *exitGameInput) PlayerIndex() int {
	return egi.playerIndex
}
func (egi *exitGameInput) ChangeState(gameObj interfaces.Game) {
	(gameObj).(*hangman).removePlayer(egi.playerIndex)
}

func (cgi *closeGameInput) GameHash() string {
	return cgi.gameHash
}
func (cgi *closeGameInput) PlayerIndex() int {
	return cgi.playerIndex
}
func (cgi *closeGameInput) ChangeState(gameObj interfaces.Game) {
}
