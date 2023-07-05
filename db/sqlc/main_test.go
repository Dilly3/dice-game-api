package db

import (
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"

	"github.com/dilly3/dice-game-api/config"
	"github.com/dilly3/dice-game-api/repository"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var C config.Configuration
var StoreIntx repository.GameRepo

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

	StoreIntx = NewPGXDB(C.DbDriverName, C.DbDataSourceName)

	os.Exit(m.Run())
}
