package main

import (
	"fmt"
	"log"
	"time"

	config "github.com/dilly3/dice-game-api/config"
	db "github.com/dilly3/dice-game-api/db/sqlc"
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
	fmt.Println("Configuration => ", config.ConfigTx.DbDataSourceName)
	fmt.Println("DRIVER =>", config.ConfigTx.DbDriverName)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func main() {

	fmt.Println("welcome to Dice Game")
	config.ConfigTx.IsGameInSession = false

	<-time.After(time.Second * 2)
	StoreIntx := db.StartDb(config.ConfigTx.DbDriverName, config.ConfigTx.DbDataSourceName)
	h := server.Setup(StoreIntx)
	s := server.NewServer(h)

	log.Fatal(s.Router.Listen(":8000"))
}
