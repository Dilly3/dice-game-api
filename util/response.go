package util

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type ResponseDto struct {
	Message       string      `json:"message,omitempty"`
	Data          interface{} `json:"data,omitempty"`
	Errors        []string    `json:"errors,omitempty"`
	Status        int         `json:"status,omitempty"`
	JackpotNumber int         `json:"jackpot_number,omitempty"`
	TimeStamp     string      `json:"timestamp,omitempty"`
}

type ResponseError struct {
	Message string `json:"message,omitempty"`
	Status  int    `json:"status,omitempty"`
}
type RollResponseDto struct {
	Message string `json:"message,omitempty"`
	Roll_1  int    `json:"roll-1,omitempty"`
	Roll_2  int    `json:"roll-2,omitempty"`
}

func RollResponse(f *fiber.Ctx, message string, roll_1 int, roll_2 int) error {
	f.SendStatus(fiber.StatusOK)
	responsedata := RollResponseDto{
		Message: message,
		Roll_1:  roll_1,
		Roll_2:  roll_2,
	}
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

func Response(f *fiber.Ctx, message string, status int, data interface{}, jackpotNum int, errs []string) error {
	f.SendStatus(status)
	responsedata := ResponseDto{
		Message:       message,
		Data:          data,
		Errors:        errs,
		JackpotNumber: jackpotNum,
		Status:        status,
		TimeStamp:     time.Now().Format("2006-01-02 15:04:05"),
	}
	return f.JSON(responsedata)

}
