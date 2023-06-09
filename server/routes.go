package server

import (
	db "github.com/dilly3/dice-game-api/db/sqlc"
	"github.com/dilly3/dice-game-api/service"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	Handler Handler
	Router  *fiber.App
}

func setup(store *db.Store) Handler {

	userServ := service.NewUserService(store)
	transServ := service.NewTransactionService(store)

	return NewHandler(transServ, userServ)
}

func NewServer(store *db.Store) Server {
	h := setup(store)
	app := fiber.New()

	v1 := app.Group("api/v1")

	v1.Get("/all", h.GetUsers())
	v1.Get("/balance/:username", h.GetWalletBalance())
	v1.Get("/credit/:username/:amount", h.CreditWallet())
	v1.Get("/debit/:username/:amount", h.DebitWallet())

	//app.Use(middleware.Timeout(60 * time.Second))

	return Server{Handler: h, Router: app}
}
