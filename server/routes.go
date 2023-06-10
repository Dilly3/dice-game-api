package server

import (
	db "github.com/dilly3/dice-game-api/db/sqlc"
	"github.com/dilly3/dice-game-api/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Server struct {
	Handler Handler
	Router  *fiber.App
}

func Setup(store db.Store) Handler {

	userServ := service.NewUserService(store)
	transServ := service.NewTransactionService(store)

	return NewHandler(transServ, userServ)
}

func NewServer(h Handler) Server {

	app := fiber.New()
	app.Use(logger.New(logger.Config{
		Format:     " ${pid} Time:${time} Status: ${status} - ${method} ${path}\n",
		TimeFormat: "02-Jan-2006",
		TimeZone:   "America/New_York",
	}))

	v1 := app.Group("api/v1")
	v1.Post("/register", h.Register())
	v1.Post("/login", h.Login())
	v1.Get("/all", h.GetUsers())
	v1.Get("/balance", h.GetWalletBalance())
	v1.Post("/credit", h.CreditWallet())
	v1.Post("/debit", h.DebitWallet())
	v1.Get("/roll", h.RollDice())
	v1.Get("/start", h.StartGame())
	v1.Get("/session", h.GetSessionState())
	v1.Get("/end", h.StopGame())
	v1.Get("/logout", h.Logout())
	v1.Get("/transactions", h.GetTransactions())

	//app.Use(middleware.Timeout(60 * time.Second))

	return Server{Handler: h, Router: app}
}
