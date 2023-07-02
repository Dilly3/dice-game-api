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

var ResetGame = func() {
	GameConfig.IsGameInSession = false
	GameConfig.LuckyNumber = 0
	GameConfig.RollNumber1 = 0
	GameConfig.RollNumber2 = 0
}
var resetRoll = func() {
	GameConfig.RollNumber1 = 0
	GameConfig.RollNumber2 = 0
}

func RollDice1(c *fiber.Ctx) error {
	num1 := rand.Int31n(6) + 1
	GameConfig.RollNumber1 = num1

	if num1 > GameConfig.LuckyNumber {
		resetRoll()
		return c.JSON(&fiber.Map{
			"Roll-1":  num1,
			"message": "you Lost, first roll is greater than jackpot number",
		})
	}

	if num1 == GameConfig.LuckyNumber {
		resetRoll()
		return c.JSON(&fiber.Map{
			"Roll-1":  num1,
			"message": "you Lost, first roll is equal to jackpot number",
		})
	}

	if GameConfig.LuckyNumber-num1 > 6 {
		resetRoll()
		return c.JSON(&fiber.Map{
			"Roll-1":  num1,
			"message": "you Lost, u need more than 6 to hit jackpot number",
		})
	}

	return c.JSON(&fiber.Map{
		"Roll-1":  num1,
		"message": fmt.Sprintf("you need %d to win", GameConfig.LuckyNumber-num1),
	})

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