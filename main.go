package main

import (
	"context"
	"fmt"
	"log"

	"os"
	"os/signal"
	"syscall"
	"time"

	config "github.com/dilly3/dice-game-api/config"
	db "github.com/dilly3/dice-game-api/db/sqlc"
	"github.com/dilly3/dice-game-api/game"
	"github.com/dilly3/dice-game-api/server"
	"github.com/dilly3/dice-game-api/service"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
)

func init() {

	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err.Error())
	}

	err = envconfig.Process("dicegame", &config.ConfigTx)

	if err != nil {
		log.Fatal(err.Error())
	}
}

func main() {

	fmt.Println("welcome to Dice Game")
	game.ResetGame()
	var err error

	db.DefaultGameRepo, err = db.NewPGXDB(config.ConfigTx.DbDriverName, config.ConfigTx.DbDataSourceName)
	<-time.After(time.Second * 2)
	if err != nil {
		log.Fatal(err)
	}

	service.DefaultGameService = service.NewGameService(db.DefaultGameRepo)
	server.FiberEngine = fiber.New()
	s := server.NewServer(config.ConfigTx.Port, server.FiberEngine)

	s.StartServer()

	errs := make(chan error, 2)

	go func() {

		errs <- s.ListenAndServe()
	}()
	c := make(chan os.Signal, 1)
	go func() {
		signal.Notify(c, syscall.SIGINT)
		ctxx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		<-ctxx.Done()

	}()

	select {
	case err := <-errs:
		log.Printf("server error: %s", err)
		os.Exit(1)
	case sig := <-c:
		log.Printf("received signal %v", sig)
		<-time.After(time.Second * 1)
		log.Printf("gracefully shutting down server...")
		<-time.After(time.Second * 1)
		os.Exit(0)

	}

}
