package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	"time"

	db "github.com/dilly3/dice-game-api/db/sqlc"
	"github.com/dilly3/dice-game-api/game"
	"github.com/dilly3/dice-game-api/service"
	"github.com/dilly3/dice-game-api/util"
	"github.com/gofiber/fiber/v2"
)

func GetUsers() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {

		users, err := service.DefaultGameService.GetAllUsers(context.Background(), db.ListUsersParams{Limit: 10, Offset: 0})

		if err != nil {
			return err
		}

		return util.Response(c, "users", fiber.StatusOK, users, 0, nil)
	}
}

func Register() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		user := &db.RegisterUserDto{}

		bytes, err := ioutil.ReadAll(bytes.NewReader(c.Body()))
		if err != nil {
			return util.ErrorResponse(c, "empty json", fiber.StatusBadRequest)
		}

		err = json.Unmarshal(bytes, user)
		if err != nil {
			return util.ErrorResponse(c, "invalid json", fiber.StatusBadRequest)
		}
		err = db.VerifyUserData(user)
		if err != nil {
			return util.ErrorResponse(c, err.Error(), fiber.StatusBadRequest)
		}
		dbuser, _ := service.DefaultGameService.GetUserByUsername(context.Background(), user.Username)

		if dbuser.ID != 0 {
			return util.ErrorResponse(c, "username already exists", fiber.StatusBadRequest)

		}

		if user.Password != user.ConfirmPassword {
			return util.ErrorResponse(c, "passwords do not match", fiber.StatusBadRequest)

		}

		dbuser, err = service.DefaultGameService.CreateUser(db.CreateUserParams{
			Firstname: user.Firstname,
			Lastname:  user.Lastname,
			Username:  user.Username,
			Password:  user.Password,
		})

		if err != nil {
			return util.ErrorResponse(c, "internal Server error :"+err.Error(), fiber.StatusBadRequest)

		}

		return util.Response(c, "user created successfully", fiber.StatusCreated, user, 0, nil)
	}
}

func Login() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		loginbody := &db.LoginDto{}

		bytes, err := ioutil.ReadAll(bytes.NewReader(c.Body()))
		if err != nil {
			return util.ErrorResponse(c, "empty json", fiber.StatusBadRequest)
		}

		err = json.Unmarshal(bytes, loginbody)
		if err != nil {
			return util.ErrorResponse(c, "bad login credentials", fiber.StatusBadRequest)
		}

		// Get first matched record
		dbuser, err := service.DefaultGameService.GetUserByUsername(context.Background(), loginbody.Username)
		if err != nil {
			return util.ErrorResponse(c, "email or password incorrect, user", fiber.StatusBadRequest)

		}

		if err := db.CompareHashAndPassword(dbuser, loginbody.Password); err != nil || dbuser.ID == 0 {
			return util.ErrorResponse(c, "email or password incorrect , pass", fiber.StatusBadRequest)

		}

		c.Cookie(&fiber.Cookie{Name: "user", Value: dbuser.Username, HTTPOnly: true, Expires: time.Now().Add(time.Hour * 24)})

		return util.Response(c, "login successful", fiber.StatusOK, nil, 0, nil)

	}
}

func GetWalletBalance() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		username := c.Cookies("user")
		if username == "" {
			return util.ErrorResponse(c, "user not logged in", fiber.StatusForbidden)

		}

		bal, _, err := service.DefaultGameService.GetWalletBalance(username)

		if err != nil && err.Error() == "cant get wallet : sql: no rows in result set" {
			c.SendStatus(fiber.StatusBadRequest)
			return util.ErrorResponse(c, "user not available", fiber.StatusBadRequest)

		}

		if err != nil {
			c.SendStatus(fiber.StatusInternalServerError)
			return util.ErrorResponse(c, "internal Server error :"+err.Error(), fiber.StatusInternalServerError)

		}

		strbal := fmt.Sprintf("%d", bal)
		return util.WalletBallanceResponse(c, strbal)

	}
}

func CreditWallet() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		username := c.Cookies("user")
		if username == "" {

			return util.ErrorResponse(c, "user not logged in", fiber.StatusForbidden)

		}

		err := service.DefaultGameService.CreditWallet(username, 155)

		if err != nil && err.Error() == "sql: no rows in result set" {

			return util.ErrorResponse(c, "user not available", fiber.StatusBadRequest)

		}

		if err != nil {
			c.SendStatus(fiber.StatusInternalServerError)
			return util.ErrorResponse(c, "internal *Server error :"+err.Error(), fiber.StatusInternalServerError)

		}

		return util.Response(c, "wallet credited", fiber.StatusOK, nil, 0, nil)
	}
}

func DebitWallet() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {

		username := c.Cookies("user")
		if username == "" {
			return util.ErrorResponse(c, "user not logged in", fiber.StatusForbidden)

		}

		body := &db.CreateWalletDto{}
		if err := c.BodyParser(body); err == fiber.ErrUnprocessableEntity {
			c.SendStatus(fiber.StatusBadRequest)
			return util.ErrorResponse(c, "bad request", fiber.StatusBadRequest)

		}

		err := service.DefaultGameService.DebitWallet(username, body.Amount)

		if err != nil && err.Error() == "sql: no rows in result set" {
			c.SendStatus(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{"message": "user not available"})

		}

		if err != nil {
			c.SendStatus(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{"message": "internal *Server error :" + err.Error()})

		}

		return c.JSON("successful")
	}
}

func GetSessionState() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {

		username := c.Cookies("user")
		if username == "" {
			c.SendStatus(fiber.StatusForbidden)
			return c.JSON(fiber.Map{"message": "user not logged in"})
		}

		return c.JSON(fiber.Map{"isSessionActive": game.GameConfig.IsGameInSession})
	}
}

func StartGame() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		user := c.Cookies("user")
		if user == "" {
			c.SendStatus(fiber.StatusForbidden)
			return c.JSON(fiber.Map{"message": "user not logged in"})
		}
		if game.GameConfig.IsGameInSession {
			return c.JSON(fiber.Map{"message": "game already in session"})
		}

		err := service.DefaultGameService.DebitWallet(user, 20)

		if err != nil && err.Error() == "sql: no rows in result set" {
			return c.JSON(fiber.Map{"message": "user not available"})

		}

		if err != nil {
			return util.ErrorResponse(c, "internal Server error :"+err.Error(), fiber.StatusInternalServerError)

		}

		game.StartGame()
		game.GameConfig.IsGameInSession = true
		jackpotNum := game.GameConfig.LuckyNumber
		return util.Response(c, "game started, debit 20 sats, roll dice. good luck!", fiber.StatusOK, nil, jackpotNum, nil)

	}
}

func RollDice() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {

		user := c.Cookies("user")
		if user == "" {
			return util.ErrorResponse(c, "user not logged in", fiber.StatusForbidden)
		}
		if !game.GameConfig.IsGameInSession {

			return util.ErrorResponse(c, "game not in session, start game first", fiber.StatusBadRequest)

		}

		if game.GameConfig.RollNumber1 == 0 {
			err := service.DefaultGameService.DebitWallet(user, 5)

			if err != nil && err.Error() == "sql: no rows in result set" {
				return util.ErrorResponse(c, "user not available", fiber.StatusBadRequest)

			}
			if err != nil {

				return util.ErrorResponse(c, "internal *Server error :"+err.Error(), fiber.StatusInternalServerError)

			}
			// roll dice 1
			game.RollDice1()
			num1 := game.GameConfig.RollNumber1
			if num1 > game.GameConfig.LuckyNumber {
				game.ResetRoll()
				return util.RollResponse(c, "you Lost, first roll is greater than jackpot number", num1, 0)

			}

			if num1 == game.GameConfig.LuckyNumber {
				game.ResetRoll()
				return util.RollResponse(c, "you Lost, first roll is equal to jackpot number", num1, 0)
			}

			if game.GameConfig.LuckyNumber-num1 > 6 {
				game.ResetRoll()
				return util.RollResponse(c, "you Lost, u need more than 6 to sit jackpot number", num1, 0)

			}

			return util.RollResponse(c, fmt.Sprintf("you need %d to win", game.GameConfig.LuckyNumber-num1), num1, 0)

		}

		// roll dice 2

		game.RollDice2()

		temp2 := game.GameConfig.RollNumber2
		temp1 := game.GameConfig.RollNumber1

		if game.GameConfig.RollNumber2 != 0 && game.GameConfig.RollNumber1+game.GameConfig.RollNumber2 == game.GameConfig.LuckyNumber {
			err := service.DefaultGameService.CreditWalletForWin(user, 10)
			if err != nil {
				return util.ErrorResponse(c, "internal server error "+err.Error(), fiber.StatusInternalServerError)

			}
			game.ResetRoll()
			return util.RollResponse(c, "WIN WIN WIN !!!!!!, You won 10 sats", temp1, temp2)
		}

		game.ResetRoll()
		return util.RollResponse(c, "you lost", temp1, temp2)

	}
}

func StopGame() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		user := c.Cookies("user")
		if user == "" {
			c.SendStatus(fiber.StatusForbidden)
			return c.JSON(fiber.Map{"message": "user not logged in"})
		}
		if !game.GameConfig.IsGameInSession {
			c.Status(fiber.StatusUnauthorized)
			return util.Response(c, "game not in session, start game first", fiber.StatusBadRequest, nil, 0, nil)

		}

		game.StopGame()
		return util.Response(c, "game stopped", fiber.StatusOK, nil, 0, nil)

	}
}

func Logout() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		user := c.Cookies("user")
		if user == "" {
			c.SendStatus(fiber.StatusForbidden)
			return util.ErrorResponse(c, "user not logged in", fiber.StatusForbidden)

		}
		c.Cookie(&fiber.Cookie{
			Value:   "",
			Name:    "user",
			Expires: time.Now().Add(-time.Hour),
		})
		game.StopGame()
		return util.Response(c, "logged out", fiber.StatusOK, nil, 0, nil)

	}
}

func GetTransactions() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		user := c.Cookies("user")

		if user == "" {
			return util.ErrorResponse(c, "user not logged in", fiber.StatusForbidden)

		}
		limit, err := strconv.Atoi(c.Params("limit"))
		if err != nil {
			return util.ErrorResponse(c, "request missing query", fiber.StatusBadRequest)
		}
		transactions, err := service.DefaultGameService.GetTransactionHistory(user, limit)
		if err != nil && err.Error() == "sql: no rows in result set" {

			return util.ErrorResponse(c, "user not available", fiber.StatusBadRequest)

		}

		if err != nil {
			return util.ErrorResponse(c, "internal *Server error :"+err.Error(), fiber.StatusInternalServerError)

		}
		return util.Response(c, "transactions", fiber.StatusOK, transactions, 0, nil)
	}
}
