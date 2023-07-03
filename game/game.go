package game

import (
	"math/rand"
)

type Game struct {
	IsGameInSession bool
	RollNumber1     int32
	RollNumber2     int32
	LuckyNumber     int32
}

var GameConfig Game

var ResetGame = func() {
	GameConfig.IsGameInSession = false
	GameConfig.LuckyNumber = 0
	GameConfig.RollNumber1 = 0
	GameConfig.RollNumber2 = 0
}

func RollDice1() {
	num1 := rand.Int31n(6) + 1
	GameConfig.RollNumber1 = num1

}

func RollDice2() {
	num2 := rand.Int31n(6) + 1
	GameConfig.RollNumber2 = num2

}

func StartGame() {

	randNum := rand.Int31n(11) + 2
	GameConfig.IsGameInSession = true
	GameConfig.LuckyNumber = randNum

}
func StopGame() {

	ResetGame()
}
