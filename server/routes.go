package server

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.uber.org/zap"
)

type Routes []Route
type handler []func(c *fiber.Ctx) error

var ServerInst *Server

type Server struct {
	Router *fiber.App
	Port   string
	Logger *zap.Logger
}

func NewServer(port string, router *fiber.App) *Server {
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("error setting up Logger =>%v\n", err)
		log.Fatal(err)
	}
	return &Server{
		Router: router,
		Port:   port,
		Logger: logger,
	}
}

func (s *Server) ListenAndServe() error {
	return s.Router.Listen(fmt.Sprintf(":%s", s.Port))
}

func (s *Server) StartServer() {
	//  start user service with default repo instance

	s.Router.Use(logger.New(logger.Config{
		Format:     " ${pid} Time:${time} Status: ${status} - ${method} ${path}\n",
		TimeFormat: "02-Jan-2006 15:04:05",
		TimeZone:   "GMT+1",
	}))

	v1 := s.Router.Group("api/v1")

	for _, route := range routes {
		v1.Add(route.Method, route.Path, route.Handler...)
	}

}

type Route struct {
	Method  string
	Path    string
	Handler handler
}

var routes = Routes{
	Route{
		Method:  "POST",
		Path:    "/register",
		Handler: handler{Register()},
	},
	Route{
		Method:  "POST",
		Path:    "/login",
		Handler: handler{Login()},
	},
	Route{
		Method:  "GET",
		Path:    "/all",
		Handler: handler{GetUsers()},
	},
	Route{
		Method:  "GET",
		Path:    "/balance",
		Handler: handler{GetWalletBalance()},
	},
	Route{
		Method:  "GET",
		Path:    "/credit",
		Handler: handler{CreditWallet()},
	},
	Route{
		Method:  "GET",
		Path:    "/roll",
		Handler: handler{RollDice()},
	},
	Route{
		Method:  "GET",
		Path:    "/start",
		Handler: handler{StartGame()},
	},
	Route{
		Method:  "GET",
		Path:    "/session",
		Handler: handler{GetSessionState()},
	},
	Route{
		Method:  "GET",
		Path:    "/end",
		Handler: handler{StopGame()},
	},
	Route{
		Method:  "GET",
		Path:    "/logout",
		Handler: handler{Logout()},
	},
	Route{
		Method:  "GET",
		Path:    "/transactions",
		Handler: handler{GetTransactions()},
	},
}
