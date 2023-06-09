package server

import (
	"context"
	"fmt"
	"log"
	"strconv"

	db "github.com/dilly3/dice-game-api/db/sqlc"
	"github.com/dilly3/dice-game-api/service"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Handler struct {
	userService *service.UserService
	trxService  *service.TransactionService
	Logger      zap.Logger
}

func NewHandler(ts *service.TransactionService, us *service.UserService) Handler {
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("error setting up Logger =>%v\n", err)
		log.Fatal(err)
	}
	return Handler{userService: us, trxService: ts, Logger: *logger}
}

func (h Handler) GetUsers() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {

		users, err := h.userService.GetAllUsers(context.Background(), db.ListUsersParams{Limit: 10, Offset: 0})

		if err != nil {
			return err
		}

		return c.JSON(users)
	}
}

func (h Handler) GetWalletBalance() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		username := c.Params("username")
		bal, err := h.userService.GetWalletBalance(username)

		if err != nil {
			return err
		}

		return c.JSON(bal)
	}
}

func (h Handler) CreditWallet() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		username := c.Params("username")
		amt, err := strconv.Atoi(c.Params("amount"))
		if err != nil {
			return fmt.Errorf("cant parse string : %v", err)
		}
		err = h.userService.CreditWallet(username, int64(amt))

		if err != nil {
			return err
		}

		return c.JSON("successful")
	}
}

func (h Handler) DebitWallet() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		username := c.Params("username")
		amt, err := strconv.Atoi(c.Params("amount"))
		if err != nil {
			return fmt.Errorf("cant parse string : %v", err)
		}
		err = h.userService.DebitWallet(username, int64(amt))

		if err != nil {
			return err
		}

		return c.JSON("successful")
	}
}
