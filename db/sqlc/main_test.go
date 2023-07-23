package db

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/dilly3/dice-game-api/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
)

func TestMain(m *testing.M) {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err.Error())
	}

	err = envconfig.Process("dicegame", &config.ConfigTx)
	fmt.Println("SPECIFICATION => ", config.ConfigTx.DbDataSourceName)
	if err != nil {
		log.Fatal(err.Error())
	}

	DefaultGameRepo, err = SetupTestDb("../../.env")
	if err != nil {
		log.Fatal(err)
	}
	TestRouter = fiber.New()

	TestRouter.Use(logger.New(logger.Config{
		Format:     " ${pid} Time:${time} Status: ${status} - ${method} ${path}\n",
		TimeFormat: "02-Jan-2006 15:04:05",
		TimeZone:   "GMT+1",
	}))
	os.Exit(m.Run())
}
