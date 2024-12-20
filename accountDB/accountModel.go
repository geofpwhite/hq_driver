package accountDB

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Account struct {
	ID       uint `gorm:"primarykey"`
	Username string
	Password string
	Wins     []*GameInstance `gorm:"many2many:account_wins"`
	Losses   []*GameInstance `gorm:"many2many:account_losses"`
}

type GameType int

const (
	HANGMAN GameType = iota
	CONNECT4
	CONNECTTHEDOTS
	TICTACTOE
)

func GenerateID(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)

}

type GameInstance struct {
	ID       uint `gorm:"primarykey"`
	GameType GameType
	Winners  []*Account `gorm:"many2many:account_wins"`
	Losers   []*Account `gorm:"many2many:account_losses"`
}

type AccountsGamesController struct {
	*sync.Mutex
	db    *gorm.DB
	users map[string]*Account
}

func NewAccountsGamesHandler() *AccountsGamesController {
	db, err := gorm.Open(sqlite.Open("./accounts.db"))
	if err != nil {
		panic("err reading accounts.db")
	}

	agh := &AccountsGamesController{
		db:    db,
		users: make(map[string]*Account),
	}
	agh.db.AutoMigrate(&Account{})
	agh.Mutex = &sync.Mutex{}
	return agh
}

func (agh *AccountsGamesController) Register(username, password string) error {
	agh.Lock()
	defer agh.Unlock()
	bytes := []byte(password)
	check := &Account{}
	agh.db.Find(check, "username='"+username+"'")
	if check.Username == username {
		return errors.New("username is taken")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword(bytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	// Comparing the password with the hash
	err = bcrypt.CompareHashAndPassword(hashedPassword, bytes) // nil means it is a match
	if err != nil {
		return err
	}

	acc := &Account{
		Username: username,
		Password: string(hashedPassword),
	}
	agh.db.Save(acc)
	return nil
}
func (agh *AccountsGamesController) Login(username, password string) (string, error) {
	agh.Lock()
	defer agh.Unlock()
	bytes := []byte(password)
	check := &Account{}
	agh.db.Find(check, "username='"+username+"'")
	err := bcrypt.CompareHashAndPassword([]byte(check.Password), bytes)
	if err != nil || agh.users[username] != nil {
		return "", errors.New("already logged in")
	}
	id := GenerateID(12)
	agh.users[id] = check
	return id, nil
}

func (agh *AccountsGamesController) Logout(id string) {
	delete(agh.users, id)
}

func (agh *AccountsGamesController) RecordGame(gameType GameType, winners, losers []*Account) {
	agh.Lock()
	defer agh.Unlock()
	game := &GameInstance{GameType: gameType, Winners: winners, Losers: losers}
	agh.db.Save(game)
}

func (agh *AccountsGamesController) AddLoser(game *GameInstance, loser *Account) {
	agh.Lock()
	defer agh.Unlock()
	game.Losers = append(game.Losers, loser)
	agh.db.Save(game)
}
func (agh *AccountsGamesController) AddWinner(game *GameInstance, winner *Account) {
	agh.Lock()
	defer agh.Unlock()
	game.Winners = append(game.Winners, winner)
	agh.db.Save(game)
}
func (agh *AccountsGamesController) GetGame(gameID uint) *GameInstance {
	agh.Lock()
	defer agh.Unlock()
	var game *GameInstance
	agh.db.Find(game, "ID="+strconv.Itoa(int(gameID)))
	if game == nil {
		return nil
	}
	return game
}

func main() {
	agh := NewAccountsGamesHandler()
	agh.Register("example", "example")
	agh.Login("example", "example")

	agh.RecordGame(HANGMAN, []*Account{agh.users["example"]}, []*Account{})

	fmt.Println(agh)
}
