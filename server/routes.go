package server

import (
	"fmt"
	"log"

	"github.com/dilly3/dice-game-api/repository"
	"github.com/dilly3/dice-game-api/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.uber.org/zap"
)

type Server struct {
	Router      *fiber.App
	Port        string
	Logger      *zap.Logger
	gameService service.GameService
}

func NewServer(port string, repo repository.GameRepo, router *fiber.App) *Server {
	gservice := service.GetGameService(repo)
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("error setting up Logger =>%v\n", err)
		log.Fatal(err)
	}
	return &Server{
		Router:      router,
		Port:        port,
		gameService: gservice,
		Logger:      logger,
	}
}

func (s *Server) Listen() error {
	return s.Router.Listen(s.Port)
}

func (s *Server) StartServer() {
	//  start user service with default repo instance

	s.Router.Use(logger.New(logger.Config{
		Format:     " ${pid} Time:${time} Status: ${status} - ${method} ${path}\n",
		TimeFormat: "02-Jan-2006 15:04:05",
		TimeZone:   "GMT+1",
	}))

	v1 := s.Router.Group("api/v1")
	v1.Post("/register", s.Register())
	v1.Post("/login", s.Login())
	v1.Get("/all", s.GetUsers())
	v1.Get("/balance", s.GetWalletBalance())
	v1.Get("/credit", s.CreditWallet())
	v1.Get("/roll", s.RollDice())
	v1.Get("/start", s.StartGame())
	v1.Get("/session", s.GetSessionState())
	v1.Get("/end", s.StopGame())
	v1.Get("/logout", s.Logout())
	v1.Get("/transactions", s.GetTransactions())

	//app.Use(middleware.Timeout(60 * time.Second))

}
