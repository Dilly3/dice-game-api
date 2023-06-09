package server

import (
	"context"
	"fmt"
	"log"

	"time"

	"github.com/dilly3/dice-game-api/config"
	db "github.com/dilly3/dice-game-api/db/sqlc"
	"github.com/dilly3/dice-game-api/service"
	"github.com/dilly3/dice-game-api/util"
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

func (h Handler) Register() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		user := &db.RegisterUserDto{}
		if err := c.BodyParser(user); err == fiber.ErrUnprocessableEntity {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{"message": "bad request"})
		}

		dbuser, _ := h.userService.GetUserByUsername(context.Background(), user.Username)
		if dbuser.ID != 0 {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{"message": "username already exists"})
		}

		if user.Password != user.ConfirmPassword {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{"message": "passwords do not match"})
		}

		dbuser, err := h.userService.CreateUser(db.CreateUserParams{
			Firstname: user.Firstname,
			Lastname:  user.Lastname,
			Username:  user.Username,
			Password:  user.Password,
		})

		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{"message": "bad request"})
		}

		c.Status(fiber.StatusCreated)
		return c.JSON(dbuser)
	}
}

func (h Handler) Login() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		loginbody := &db.LoginDto{}
		if err := c.BodyParser(loginbody); err == fiber.ErrUnprocessableEntity {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{"message": "bad login credentials"})
		}
		dbuser := db.User{}
		// Get first matched record
		dbuser, err := h.userService.GetUserByUsername(context.Background(), loginbody.Username)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{"message": "email or password incorrect"})
		}

		if err := dbuser.CompareHashAndPassword(loginbody.Password); err != nil || dbuser.ID == 0 {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{"message": "email or password incorrect"})
		}

		c.Cookie(&fiber.Cookie{Name: "user", Value: dbuser.Username, HTTPOnly: true, Expires: time.Now().Add(time.Hour * 24)})
		c.Status(fiber.StatusOK)

		return util.Response(c, "login successful", fiber.StatusOK, nil, nil)

	}
}

func (h Handler) GetWalletBalance() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		username := c.Params("username")
		bal, assts, err := h.userService.GetWalletBalance(username)

		if err != nil {
			return err
		}

		return c.JSON(fiber.Map{"balance": bal, "assets": assts})
	}
}

func (h Handler) CreditWallet() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		body := &db.CreateWalletDto{}
		if err := c.BodyParser(body); err == fiber.ErrUnprocessableEntity {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{"message": "bad request"})
		}
		username := c.Cookies("user")

		err := h.userService.CreditWallet(username, int32(body.Amount))

		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{"message": err.Error()})

		}

		return c.JSON("successful")
	}
}

func (h Handler) DebitWallet() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		body := &db.CreateWalletDto{}
		if err := c.BodyParser(body); err == fiber.ErrUnprocessableEntity {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{"message": "bad request"})
		}

		username := c.Cookies("user")

		err := h.userService.DebitWallet(username, int32(body.Amount))

		if err != nil {
			return err
		}

		return c.JSON("successful")
	}
}

func (h Handler) GetSessionState() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {

		return c.JSON(fiber.Map{"isSessionActive": config.ConfigTx.IsGameInSession})
	}
}

func (h Handler) StartGame() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		user := c.Cookies("user")
		if user == "" {
			return c.JSON(fiber.Map{"message": "user not logged in"})
		}
		if config.ConfigTx.IsGameInSession {
			return c.JSON(fiber.Map{"message": "game already in session"})
		}

		err := h.userService.DebitWallet(user, 20)

		if err != nil {
			return c.JSON(fiber.Map{"message": err.Error()})

		}

		config.StartGame()
		config.ConfigTx.IsGameInSession = true
		config.ConfigTx.NumberOfTrials = 10

		c.JSON(fiber.Map{"message": "game started, roll dice. good luck!", "debit": "20 sats", "luckyNumber": config.ConfigTx.LuckyNumber, "isSessionActive": config.ConfigTx.IsGameInSession})

		return nil

	}
}

func (h Handler) RollDice() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {

		user := c.Cookies("user")
		if user == "" {
			return c.JSON(fiber.Map{"message": "user not logged in"})
		}
		if !config.ConfigTx.IsGameInSession {

			return c.JSON(fiber.Map{"message": "game not in session, start game first"})
		}

		err := h.userService.DebitWallet(user, 5)

		if err != nil {
			return c.JSON(fiber.Map{"message": err.Error()})

		}
		// roll dice
		res, err := config.RollDice()

		if err != nil {

			c.JSON(fiber.Map{"message": err.Error()})
		}
		if res.RollNumber1+res.RollNumber2 == config.ConfigTx.LuckyNumber {
			err := h.userService.CreditWallet(user, 10)
			if err != nil {
				return c.JSON(fiber.Map{"message": err.Error()})
			}
			config.ConfigTx.LuckyNumber = 0
			config.ConfigTx.RollNumber1 = 0
			config.ConfigTx.RollNumber2 = 0

			return c.JSON(fiber.Map{"message": "you won 10 credits", "result": res})
		}
		return c.JSON(fiber.Map{"message": "you lost", "result": res})

	}

}

func (h Handler) StopGame() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		user := c.Cookies("user")
		if user == "" {
			return c.JSON(fiber.Map{"message": "user not logged in"})
		}
		if !config.ConfigTx.IsGameInSession {
			return c.JSON(fiber.Map{"message": "game not in session, start game first"})
		}

		config.ConfigTx.IsGameInSession = false
		config.ConfigTx.NumberOfTrials = 0
		config.ConfigTx.LuckyNumber = 0
		config.ConfigTx.RollNumber1 = 0
		config.ConfigTx.RollNumber2 = 0

		return c.JSON(fiber.Map{"message": "game stopped"})
	}
}

func (h Handler) Logout() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		user := c.Cookies("user")
		if user == "" {
			return c.JSON(fiber.Map{"message": "user not logged in"})
		}
		c.Cookie(&fiber.Cookie{
			Value:   "",
			Name:    "user",
			Expires: time.Now().Add(-time.Hour),
		})

		config.ConfigTx.IsGameInSession = false
		config.ConfigTx.NumberOfTrials = 0
		config.ConfigTx.LuckyNumber = 0
		config.ConfigTx.RollNumber1 = 0
		config.ConfigTx.RollNumber2 = 0

		return c.JSON(fiber.Map{"message": "logged out"})
	}
}

func (h Handler) GetTransactions() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		user := c.Cookies("user")
		if user == "" {
			return c.JSON(fiber.Map{"message": "user not logged in"})
		}
		transactions, err := h.trxService.GetTransactionHistory(user)
		if err != nil {
			return c.JSON(fiber.Map{"message": err.Error()})
		}
		return c.JSON(fiber.Map{"transactions": transactions})
	}
}
