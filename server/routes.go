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
		response := util.NewResponseBuilder().SetMessage("users").SetStatus(fiber.StatusOK).SetData(users).Build()
		return util.Response(c, response)
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

		dbuser, err = service.DefaultGameService.CreateUser(*user)

		if err != nil {
			return util.ErrorResponse(c, "internal Server error :"+err.Error(), fiber.StatusBadRequest)

		}
		response := util.NewResponseBuilder().SetMessage("user created successfully").SetStatus(fiber.StatusCreated).SetData(dbuser).Build()
		return util.Response(c, response)
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

		if loginbody.Username == "" || loginbody.Password == "" {
			return util.ErrorResponse(c, "bad login credentials", fiber.StatusBadRequest)
		}

		// Get first matched record
		dbuser, err := service.DefaultGameService.GetUserByUsername(context.Background(), loginbody.Username)
		if err != nil {
			return util.ErrorResponse(c, "username or password incorrect", fiber.StatusBadRequest)

		}

		if err := db.CompareHashAndPassword(dbuser, loginbody.Password); err != nil || dbuser.ID == 0 {
			return util.ErrorResponse(c, "username or password incorrect , pass", fiber.StatusBadRequest)

		}

		c.Cookie(&fiber.Cookie{Name: "user", Value: dbuser.Username, HTTPOnly: true, Expires: time.Now().Add(time.Hour * 24)})

		response := util.NewResponseBuilder().SetMessage("login successful").SetStatus(fiber.StatusOK).SetTimeStamp().Build()

		return util.Response(c, response)

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
		response := util.NewResponseBuilder().SetStatus(fiber.StatusOK).SetAssets().SetBalance(strbal).
			SetTimeStamp().
			Build()
		return util.Response(c, response)

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
			return util.ErrorResponse(c, err.Error(), fiber.StatusInternalServerError)

		}
		response := util.NewResponseBuilder().SetMessage("wallet credited").SetStatus(fiber.StatusOK).Build()
		return util.Response(c, response)
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
		response := util.NewResponseBuilder().SetStatus(fiber.StatusOK).SetSessionState(game.GameConfig.IsGameInSession).SetTimeStamp().Build()
		return util.Response(c, response)
	}
}

func StartGame() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		user := c.Cookies("user")
		if user == "" {
			c.SendStatus(fiber.StatusForbidden)
			return c.JSON(fiber.Map{"message": "user not logged in"})
		}
		// if game.GameConfig.IsGameInSession {
		// 	return c.JSON(fiber.Map{"message": "game already in session"})
		// }

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
		response := util.NewResponseBuilder().SetMessage("game started, debit 20 sats, roll dice. good luck!").
			SetStatus(fiber.StatusOK).
			SetJackpotNumber(jackpotNum).
			Build()
		return util.Response(c, response)

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

		if game.GameConfig.Gamescore.RollNumber1 == 0 {
			err := service.DefaultGameService.DebitWallet(user, 5)

			if err != nil && err.Error() == "sql: no rows in result set" {
				return util.ErrorResponse(c, "user not available", fiber.StatusBadRequest)

			}
			if err != nil {

				return util.ErrorResponse(c, "internal *Server error :"+err.Error(), fiber.StatusInternalServerError)

			}
			// roll dice 1
			game.RollDice1()
			num1 := game.GameConfig.Gamescore.RollNumber1
			if num1 > game.GameConfig.LuckyNumber {
				game.ResetRoll()
				rollResponse := util.NewResponseBuilder().SetMessage("you Lost, first roll is greater than jackpot number").SetGameScore(num1, 0).
					SetJackpotNumber(game.GameConfig.LuckyNumber).Build()
				return util.Response(c, rollResponse)

			}

			if num1 == game.GameConfig.LuckyNumber {
				game.ResetRoll()
				rollResponse := util.NewResponseBuilder().
					SetMessage("you lost, first roll is equal to jackpot number").
					SetGameScore(num1, 0).
					SetJackpotNumber(game.GameConfig.LuckyNumber).Build()
				return util.Response(c, rollResponse)
			}

			if game.GameConfig.LuckyNumber-num1 > 6 {
				game.ResetRoll()
				rollResponse := util.NewResponseBuilder().SetMessage("you Lost, u need more than 6 to sit jackpot number").SetGameScore(num1, 0).
					SetJackpotNumber(game.GameConfig.LuckyNumber).Build()
				return util.Response(c, rollResponse)

			}
			rollResponse := util.NewResponseBuilder().SetMessage(fmt.Sprintf("you need %d to win", game.GameConfig.LuckyNumber-num1)).SetGameScore(num1, 0).
				SetJackpotNumber(game.GameConfig.LuckyNumber).
				Build()
			return util.Response(c, rollResponse)

		}

		// roll dice 2

		game.RollDice2()

		temp2 := game.GameConfig.Gamescore.RollNumber2
		temp1 := game.GameConfig.Gamescore.RollNumber1

		if game.GameConfig.Gamescore.RollNumber2 != 0 && game.GameConfig.Gamescore.RollNumber1+game.GameConfig.Gamescore.RollNumber2 == game.GameConfig.LuckyNumber {
			err := service.DefaultGameService.CreditWalletForWin(user, 10)
			if err != nil {
				return util.ErrorResponse(c, "internal server error "+err.Error(), fiber.StatusInternalServerError)

			}
			game.ResetRoll()
			rollResponse := util.NewResponseBuilder().SetMessage("WIN WIN WIN !!!!!!, You won 10 sats").SetGameScore(temp1, temp2).
				SetJackpotNumber(game.GameConfig.LuckyNumber).
				Build()
			return util.Response(c, rollResponse)
		}

		game.ResetRoll()
		rollResponse := util.NewResponseBuilder().SetMessage("you Lost, try again").SetGameScore(temp1, temp2).
			SetJackpotNumber(game.GameConfig.LuckyNumber).
			Build()
		return util.Response(c, rollResponse)

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
			response := util.NewResponseBuilder().SetMessage("game not in session, start game first").
				SetStatus(fiber.StatusBadRequest).
				Build()

			return util.Response(c, response)

		}

		game.StopGame()
		response := util.NewResponseBuilder().SetMessage("game stopped").
			SetStatus(fiber.StatusOK).
			Build()
		return util.Response(c, response)

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
		response := util.NewResponseBuilder().SetMessage("user logged out").
			SetStatus(fiber.StatusOK).
			Build()

		return util.Response(c, response)

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
		response := util.NewResponseBuilder().SetMessage("transactions").
			SetStatus(fiber.StatusOK).
			SetData(transactions).
			Build()

		return util.Response(c, response)
	}
}
