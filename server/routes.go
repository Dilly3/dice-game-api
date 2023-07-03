package server

import (
	"github.com/dilly3/dice-game-api/repository"
	"github.com/dilly3/dice-game-api/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Server struct {
	Router *fiber.App
}

func StartServer(repo repository.GameRepo) Server {
	//  start user service with default repo instance
	userserv := service.NewGameService(repo)
	h := NewHandler(*userserv)
	app := fiber.New()
	app.Use(logger.New(logger.Config{
		Format:     " ${pid} Time:${time} Status: ${status} - ${method} ${path}\n",
		TimeFormat: "02-Jan-2006 15:04:05",
		TimeZone:   "GMT+1",
	}))

	v1 := app.Group("api/v1")
	v1.Post("/register", h.Register())
	v1.Post("/login", h.Login())
	v1.Get("/all", h.GetUsers())
	v1.Get("/balance", h.GetWalletBalance())
	v1.Get("/credit", h.CreditWallet())
	v1.Get("/roll", h.RollDice())
	v1.Get("/start", h.StartGame())
	v1.Get("/session", h.GetSessionState())
	v1.Get("/end", h.StopGame())
	v1.Get("/logout", h.Logout())
	v1.Get("/transactions", h.GetTransactions())

	//app.Use(middleware.Timeout(60 * time.Second))

	return Server{Router: app}
}
