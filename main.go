package main

import (
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
	"github.com/joho/godotenv"

	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
)

var err error
var StoreIntx db.Store

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
	game.GameConfig.IsGameInSession = false

	<-time.After(time.Second * 2)
	db.StartDb(config.ConfigTx.DbDriverName, config.ConfigTx.DbDataSourceName)

	s := server.StartServer()

	errs := make(chan error, 2)

	go func() {

		errs <- s.Router.Listen(":8000")
	}()
	c := make(chan os.Signal, 1)
	go func() {
		signal.Notify(c, syscall.SIGINT)

	}()

	select {
	case err := <-errs:
		log.Printf("server error: %s", err)
		os.Exit(1)
	case sig := <-c:
		log.Printf("received signal %s", sig)
		<-time.After(time.Second * 1)
		log.Printf("gracefully shutting down server...")
		<-time.After(time.Second * 1)
		os.Exit(0)
	}
}
