package util

import (
	"github.com/gofiber/fiber/v2"
)

type ResponseDto struct {
	Status    int      `json:"status,omitempty"`
	GameScore any      `json:"gameScore,omitempty"`
	Message   string   `json:"message,omitempty"`
	Data      any      `json:"data,omitempty"`
	Errors    []string `json:"errors,omitempty"`

	Balance       string `json:"balance,omitempty"`
	Assts         string `json:"assets,omitempty"`
	JackpotNumber int    `json:"jackpot_number,omitempty"`
	TimeStamp     string `json:"timestamp,omitempty"`
	SessionState  bool   `json:"session_state,omitempty"`
}

type ResponseError struct {
	Message string `json:"message,omitempty"`
	Status  int    `json:"status,omitempty"`
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

	return f.JSON(args)

}
