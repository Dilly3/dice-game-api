package server

import (
	"context"
	"fmt"
	"log"

	"time"

	"github.com/dilly3/dice-game-api/game"
	"github.com/dilly3/dice-game-api/models"
	"github.com/dilly3/dice-game-api/service"
	"github.com/dilly3/dice-game-api/util"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Handler struct {
	Logger      zap.Logger
	gameService service.GameService
}

var resetRoll = func() {
	game.GameConfig.RollNumber1 = 0
	game.GameConfig.RollNumber2 = 0
}

func NewHandler(userser service.GameService) Handler {

	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("error setting up Logger =>%v\n", err)
		log.Fatal(err)
	}
	return Handler{Logger: *logger, gameService: userser}
}

func (h Handler) GetUsers() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {

		users, err := h.gameService.GetAllUsers(context.Background(), models.ListUsersParams{Limit: 10, Offset: 0})

		if err != nil {
			return err
		}

		return c.JSON(users)
	}
}

func (h Handler) Register() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		user := &models.RegisterUserDto{}
		if err := c.BodyParser(user); err == fiber.ErrUnprocessableEntity {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{"message": "bad request"})
		}

		dbuser, _ := h.gameService.GetUserByUsername(context.Background(), user.Username)
		if dbuser.ID != 0 {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{"message": "username already exists"})
		}

		if user.Password != user.ConfirmPassword {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{"message": "passwords do not match"})
		}

		dbuser, err := h.gameService.CreateUser(models.CreateUserParams{
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
		loginbody := &models.LoginDto{}
		if err := c.BodyParser(loginbody); err == fiber.ErrUnprocessableEntity {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{"message": "bad login credentials"})
		}
		dbuser := models.User{}
		// Get first matched record
		dbuser, err := h.gameService.GetUserByUsername(context.Background(), loginbody.Username)
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
		username := c.Cookies("user")
		if username == "" {
			c.SendStatus(fiber.StatusForbidden)
			return c.JSON(fiber.Map{"message": "user not logged in"})
		}

		bal, assts, err := h.gameService.GetWalletBalance(username)

		if err != nil && err.Error() == "cant get wallet : sql: no rows in result set" {
			c.SendStatus(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{"message": "user not available"})

		}

		if err != nil {
			c.SendStatus(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{"message": "internal server error :" + err.Error()})

		}

		strbal := fmt.Sprintf("%d", bal)

		return c.JSON(fiber.Map{"balance": strbal, "assets": assts})
	}
}

func (h Handler) CreditWallet() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		username := c.Cookies("user")
		if username == "" {
			c.SendStatus(fiber.StatusForbidden)
			return c.JSON(fiber.Map{"message": "user not logged in"})
		}

		err := h.gameService.CreditWallet(username, 155)

		if err != nil && err.Error() == "sql: no rows in result set" {
			c.SendStatus(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{"message": "user not available"})

		}

		if err != nil {
			c.SendStatus(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{"message": "internal server error :" + err.Error()})

		}

		return c.JSON("successful")
	}
}

func (h Handler) DebitWallet() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {

		username := c.Cookies("user")
		if username == "" {
			c.SendStatus(fiber.StatusForbidden)
			return c.JSON(fiber.Map{"message": "user not logged in"})
		}

		body := &models.CreateWalletDto{}
		if err := c.BodyParser(body); err == fiber.ErrUnprocessableEntity {
			c.SendStatus(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{"message": "bad request"})
		}

		err := h.gameService.DebitWallet(username, int32(body.Amount))

		if err != nil && err.Error() == "sql: no rows in result set" {
			c.SendStatus(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{"message": "user not available"})

		}

		if err != nil {
			c.SendStatus(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{"message": "internal server error :" + err.Error()})

		}

		return c.JSON("successful")
	}
}

func (h Handler) GetSessionState() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {

		username := c.Cookies("user")
		if username == "" {
			c.SendStatus(fiber.StatusForbidden)
			return c.JSON(fiber.Map{"message": "user not logged in"})
		}

		return c.JSON(fiber.Map{"isSessionActive": game.GameConfig.IsGameInSession})
	}
}

func (h Handler) StartGame() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		user := c.Cookies("user")
		if user == "" {
			c.SendStatus(fiber.StatusForbidden)
			return c.JSON(fiber.Map{"message": "user not logged in"})
		}
		if game.GameConfig.IsGameInSession {
			return c.JSON(fiber.Map{"message": "game already in session"})
		}

		err := h.gameService.DebitWallet(user, 20)

		if err != nil && err.Error() == "sql: no rows in result set" {
			c.SendStatus(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{"message": "user not available"})

		}

		if err != nil {
			c.SendStatus(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{"message": "internal server error :" + err.Error()})

		}

		game.StartGame()
		game.GameConfig.IsGameInSession = true

		c.JSON(fiber.Map{"message": "game started, roll dice. good luck!", "debit": "20 sats", "JackpotNumber": game.GameConfig.LuckyNumber, "isSessionActive": game.GameConfig.IsGameInSession})

		return nil

	}
}

func (h Handler) RollDice() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {

		user := c.Cookies("user")
		if user == "" {
			c.SendStatus(fiber.StatusForbidden)
			return c.JSON(fiber.Map{"message": "user not logged in"})
		}
		if !game.GameConfig.IsGameInSession {
			c.SendStatus(fiber.StatusBadRequest)

			return c.JSON(fiber.Map{"message": "game not in session, start game first"})
		}

		if game.GameConfig.RollNumber1 == 0 {
			err := h.gameService.DebitWallet(user, 5)

			if err != nil && err.Error() == "sql: no rows in result set" {
				c.SendStatus(fiber.StatusBadRequest)
				return c.JSON(fiber.Map{"message": "user not available"})

			}
			if err != nil {
				c.SendStatus(fiber.StatusInternalServerError)
				return c.JSON(fiber.Map{"message": "internal server error :" + err.Error()})

			}
			// roll dice 1
			game.RollDice1()
			num1 := game.GameConfig.RollNumber1
			if num1 > game.GameConfig.LuckyNumber {
				resetRoll()
				return c.JSON(&fiber.Map{
					"Roll-1":  num1,
					"message": "you Lost, first roll is greater than jackpot number",
				})
			}

			if num1 == game.GameConfig.LuckyNumber {
				resetRoll()
				return c.JSON(&fiber.Map{
					"Roll-1":  num1,
					"message": "you Lost, first roll is equal to jackpot number",
				})
			}

			if game.GameConfig.LuckyNumber-num1 > 6 {
				resetRoll()
				return c.JSON(&fiber.Map{
					"Roll-1":  num1,
					"message": "you Lost, u need more than 6 to hit jackpot number",
				})
			}

			return c.JSON(&fiber.Map{
				"Roll-1":  num1,
				"message": fmt.Sprintf("you need %d to win", game.GameConfig.LuckyNumber-num1),
			})
		}

		// roll dice 2

		game.RollDice2()

		temp2 := game.GameConfig.RollNumber2

		if game.GameConfig.RollNumber2 != 0 && game.GameConfig.RollNumber1+game.GameConfig.RollNumber2 == game.GameConfig.LuckyNumber {
			err := h.gameService.CreditWalletForWin(user, 10)
			if err != nil {
				return c.JSON(fiber.Map{"message": err.Error()})
			}
			resetRoll()
			return c.JSON(fiber.Map{"message": "WIN WIN WIN !!!!!!, You won 10 sats", "Roll-2": temp2})
		}

		resetRoll()
		return c.JSON(fiber.Map{"message": "you lost", "Roll-2": temp2})

	}
}

func (h Handler) StopGame() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		user := c.Cookies("user")
		if user == "" {
			c.SendStatus(fiber.StatusForbidden)
			return c.JSON(fiber.Map{"message": "user not logged in"})
		}
		if !game.GameConfig.IsGameInSession {
			c.Status(fiber.StatusUnauthorized)
			return c.JSON(fiber.Map{"message": "game not in session, start game first"})
		}

		game.StopGame()
		return c.JSON(fiber.Map{"message": "game stopped"})
	}
}

func (h Handler) Logout() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		user := c.Cookies("user")
		if user == "" {
			c.SendStatus(fiber.StatusForbidden)
			return c.JSON(fiber.Map{"message": "user not logged in"})
		}
		c.Cookie(&fiber.Cookie{
			Value:   "",
			Name:    "user",
			Expires: time.Now().Add(-time.Hour),
		})
		game.StopGame()

		return c.JSON(fiber.Map{"message": "logged out"})
	}
}

func (h Handler) GetTransactions() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		user := c.Cookies("user")
		if user == "" {
			c.SendStatus(fiber.StatusForbidden)
			return c.JSON(fiber.Map{"message": "user not logged in"})
		}
		transactions, err := h.gameService.GetTransactionHistory(user)
		if err != nil && err.Error() == "sql: no rows in result set" {
			c.SendStatus(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{"message": "user not available"})

		}

		if err != nil {
			c.SendStatus(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{"message": "internal server error" + err.Error()})

		}
		return c.JSON(fiber.Map{"transactions": transactions})
	}
}
