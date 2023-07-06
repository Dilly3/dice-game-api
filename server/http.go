package server

import (
	"context"
	"fmt"

	"time"

	"github.com/dilly3/dice-game-api/game"
	"github.com/dilly3/dice-game-api/models"
	"github.com/dilly3/dice-game-api/util"
	"github.com/gofiber/fiber/v2"
)

func (s *Server) GetUsers() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {

		users, err := s.gameService.GetAllUsers(context.Background(), models.ListUsersParams{Limit: 10, Offset: 0})

		if err != nil {
			return err
		}

		return c.JSON(users)
	}
}

func (s *Server) Register() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		user := &models.RegisterUserDto{}
		if err := c.BodyParser(user); err == fiber.ErrUnprocessableEntity {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{"message": "bad request"})
		}

		dbuser, _ := s.gameService.GetUserByUsername(context.Background(), user.Username)
		if dbuser.ID != 0 {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{"message": "username already exists"})
		}

		if user.Password != user.ConfirmPassword {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{"message": "passwords do not matcs"})
		}

		dbuser, err := s.gameService.CreateUser(models.CreateUserParams{
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

func (s *Server) Login() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		loginbody := &models.LoginDto{}
		if err := c.BodyParser(loginbody); err == fiber.ErrUnprocessableEntity {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{"message": "bad login credentials"})
		}
		dbuser := models.User{}
		// Get first matcsed record
		dbuser, err := s.gameService.GetUserByUsername(context.Background(), loginbody.Username)
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

func (s *Server) GetWalletBalance() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		username := c.Cookies("user")
		if username == "" {
			c.SendStatus(fiber.StatusForbidden)
			return c.JSON(fiber.Map{"message": "user not logged in"})
		}

		bal, assts, err := s.gameService.GetWalletBalance(username)

		if err != nil && err.Error() == "cant get wallet : sql: no rows in result set" {
			c.SendStatus(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{"message": "user not available"})

		}

		if err != nil {
			c.SendStatus(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{"message": "internal *Server error :" + err.Error()})

		}

		strbal := fmt.Sprintf("%d", bal)

		return c.JSON(fiber.Map{"balance": strbal, "assets": assts})
	}
}

func (s *Server) CreditWallet() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		username := c.Cookies("user")
		if username == "" {
			c.SendStatus(fiber.StatusForbidden)
			return c.JSON(fiber.Map{"message": "user not logged in"})
		}

		err := s.gameService.CreditWallet(username, 155)

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

func (s *Server) DebitWallet() func(*fiber.Ctx) error {
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

		err := s.gameService.DebitWallet(username, int32(body.Amount))

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

func (s *Server) GetSessionState() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {

		username := c.Cookies("user")
		if username == "" {
			c.SendStatus(fiber.StatusForbidden)
			return c.JSON(fiber.Map{"message": "user not logged in"})
		}

		return c.JSON(fiber.Map{"isSessionActive": game.GameConfig.IsGameInSession})
	}
}

func (s *Server) StartGame() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		user := c.Cookies("user")
		if user == "" {
			c.SendStatus(fiber.StatusForbidden)
			return c.JSON(fiber.Map{"message": "user not logged in"})
		}
		if game.GameConfig.IsGameInSession {
			return c.JSON(fiber.Map{"message": "game already in session"})
		}

		err := s.gameService.DebitWallet(user, 20)

		if err != nil && err.Error() == "sql: no rows in result set" {
			c.SendStatus(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{"message": "user not available"})

		}

		if err != nil {
			c.SendStatus(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{"message": "internal *Server error :" + err.Error()})

		}

		game.StartGame()
		game.GameConfig.IsGameInSession = true

		c.JSON(fiber.Map{"message": "game started, roll dice. good luck!", "debit": "20 sats", "JackpotNumber": game.GameConfig.LuckyNumber, "isSessionActive": game.GameConfig.IsGameInSession})

		return nil

	}
}

func (s *Server) RollDice() func(*fiber.Ctx) error {
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
			err := s.gameService.DebitWallet(user, 5)

			if err != nil && err.Error() == "sql: no rows in result set" {
				c.SendStatus(fiber.StatusBadRequest)
				return c.JSON(fiber.Map{"message": "user not available"})

			}
			if err != nil {
				c.SendStatus(fiber.StatusInternalServerError)
				return c.JSON(fiber.Map{"message": "internal *Server error :" + err.Error()})

			}
			// roll dice 1
			game.RollDice1()
			num1 := game.GameConfig.RollNumber1
			if num1 > game.GameConfig.LuckyNumber {
				game.ResetRoll()
				return c.JSON(&fiber.Map{
					"Roll-1":  num1,
					"message": "you Lost, first roll is greater tsan jackpot number",
				})
			}

			if num1 == game.GameConfig.LuckyNumber {
				game.ResetRoll()
				return c.JSON(&fiber.Map{
					"Roll-1":  num1,
					"message": "you Lost, first roll is equal to jackpot number",
				})
			}

			if game.GameConfig.LuckyNumber-num1 > 6 {
				game.ResetRoll()
				return c.JSON(&fiber.Map{
					"Roll-1":  num1,
					"message": "you Lost, u need more tsan 6 to sit jackpot number",
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
			err := s.gameService.CreditWalletForWin(user, 10)
			if err != nil {
				return c.JSON(fiber.Map{"message": err.Error()})
			}
			game.ResetRoll()
			return c.JSON(fiber.Map{"message": "WIN WIN WIN !!!!!!, You won 10 sats", "Roll-2": temp2})
		}

		game.ResetRoll()
		return c.JSON(fiber.Map{"message": "you lost", "Roll-2": temp2})

	}
}

func (s *Server) StopGame() func(*fiber.Ctx) error {
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

func (s *Server) Logout() func(*fiber.Ctx) error {
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

func (s *Server) GetTransactions() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		user := c.Cookies("user")
		if user == "" {
			c.SendStatus(fiber.StatusForbidden)
			return c.JSON(fiber.Map{"message": "user not logged in"})
		}
		transactions, err := s.gameService.GetTransactionHistory(user)
		if err != nil && err.Error() == "sql: no rows in result set" {
			c.SendStatus(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{"message": "user not available"})

		}

		if err != nil {
			c.SendStatus(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{"message": "internal *Server error" + err.Error()})

		}
		return c.JSON(fiber.Map{"transactions": transactions})
	}
}
