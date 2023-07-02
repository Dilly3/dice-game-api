package game

import (
	"fmt"
	"math/rand"

	"github.com/gofiber/fiber/v2"
)

type Game struct {
	IsGameInSession bool
	RollNumber1     int32
	RollNumber2     int32
	LuckyNumber     int32
}

var GameConfig Game

func RollDice1(c *fiber.Ctx) error {
	num1 := rand.Int31n(6) + 1
	GameConfig.RollNumber1 = num1
	c.JSON(&fiber.Map{
		"Roll-1":  num1,
		"message": fmt.Sprintf("you need %d to win", GameConfig.LuckyNumber-num1),
	})
	return nil
}

func RollDice2() int32 {
	num2 := rand.Int31n(6) + 1
	GameConfig.RollNumber2 = num2

	return num2
}

func StartGame() {

	randNum := rand.Int31n(11) + 2
	GameConfig.IsGameInSession = true
	GameConfig.LuckyNumber = randNum

}
func StopGame() {

	GameConfig.IsGameInSession = false
	GameConfig.LuckyNumber = 0
	GameConfig.RollNumber1 = 0
	GameConfig.RollNumber2 = 0
}
