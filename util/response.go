package util

import (
	"time"

	"github.com/dilly3/dice-game-api/game"
	"github.com/gofiber/fiber/v2"
)

type ResponseDto struct {
	Message       string      `json:"message,omitempty"`
	Data          interface{} `json:"data,omitempty"`
	Errors        []string    `json:"errors,omitempty"`
	Status        int         `json:"status,omitempty"`
	Balance       string      `json:"balance,omitempty"`
	Assts         string      `json:"assets,omitempty"`
	JackpotNumber int         `json:"jackpot_number,omitempty"`
	TimeStamp     string      `json:"timestamp,omitempty"`
}

type ResponseError struct {
	Message string `json:"message,omitempty"`
	Status  int    `json:"status,omitempty"`
}
type RollResponseDto struct {
	Message   string         `json:"message,omitempty"`
	Gamescore game.GameScore `json:"game_score,omitempty"`
	JackpotNumber int         `json:"jackpot_number,omitempty"`
}

func RollResponse(f *fiber.Ctx, arg RollResponseDto) error {
	f.SendStatus(fiber.StatusOK)
	responsedata := arg

	return f.JSON(responsedata)
}

type walletBallanceResponseDto struct {
	Balance string `json:"balance,omitempty"`
	Assts   string `json:"assets,omitempty"`
}

func WalletBallanceResponse(f *fiber.Ctx, balance string) error {
	f.SendStatus(fiber.StatusOK)
	responsedata := walletBallanceResponseDto{
		Balance: balance,
		Assts:   "sats",
	}
	return f.JSON(responsedata)
}
func ErrorResponse(f *fiber.Ctx, message string, status int) error {
	f.SendStatus(status)
	responsedata := ResponseError{
		Message: message,
		Status:  status,
	}
	return f.JSON(responsedata)
}

func Response(f *fiber.Ctx, args *ResponseDto) error {
	f.SendStatus(args.Status)
	args.TimeStamp = time.Now().Format("2006-01-02 15:04:05")

	return f.JSON(args)

}
