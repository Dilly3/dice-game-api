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

	app.Get("/all", h.GetUsers())

	//app.Use(middleware.Timeout(60 * time.Second))

	return Server{Handler: h, Router: app}
}
