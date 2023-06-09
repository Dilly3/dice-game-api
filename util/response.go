package util

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type ResponseDto struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Errors  []string    `json:"errors"`
	Status  int         `json:"status"`

	TimeStamp string `json:"timestamp"`
}

func Response(f *fiber.Ctx, message string, status int, data interface{}, errs []string) error {
	responsedata := ResponseDto{
		Message:   message,
		Data:      data,
		Errors:    errs,
		Status:    status,
		TimeStamp: time.Now().Format("2006-01-02 15:04:05"),
	}
	return f.JSON(responsedata)

}
