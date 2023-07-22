package game

import (
	"math/rand"
)

type Game struct {
	IsGameInSession bool
	Gamescore       GameScore
	LuckyNumber     int
}

type GameScore struct {
	RollNumber1 int `json:"roll1,omitempty"`
	RollNumber2 int `json:"roll2,omitempty"`
}

var GameConfig Game

var ResetGame = func() {
	GameConfig.IsGameInSession = false
	GameConfig.LuckyNumber = 0
	GameConfig.Gamescore.RollNumber1 = 0
	GameConfig.Gamescore.RollNumber2 = 0
}
var ResetRoll = func() {
	GameConfig.Gamescore.RollNumber1 = 0
	GameConfig.Gamescore.RollNumber2 = 0
}

func RollDice1() {
	num1 := rand.Intn(6) + 1
	GameConfig.Gamescore.RollNumber1 = num1

}

func RollDice2() {
	num2 := rand.Intn(6) + 1
	GameConfig.Gamescore.RollNumber2 = num2

}

func StartGame() {

	randNum := rand.Intn(11) + 2
	GameConfig.IsGameInSession = true
	GameConfig.LuckyNumber = randNum

}
func StopGame() {

	ResetGame()
}
