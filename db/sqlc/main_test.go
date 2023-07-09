package db

import (
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"

	"github.com/dilly3/dice-game-api/config"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var C config.Configuration
var StoreIntx IGameRepo
var err error

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err.Error())
	}

	err = envconfig.Process("dicegame", &C)
	fmt.Println("SPECIFICATION => ", C.DbDataSourceName)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func TestMain(m *testing.M) {

	StoreIntx, err = NewPGXDB(C.DbDriverName, C.DbDataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}
