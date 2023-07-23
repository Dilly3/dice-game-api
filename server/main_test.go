package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/dilly3/dice-game-api/config"
	db "github.com/dilly3/dice-game-api/db/sqlc"
	"github.com/dilly3/dice-game-api/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

func TestMain(m *testing.M) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err.Error())
	}

	err = envconfig.Process("dicegame", &config.ConfigTx)
	fmt.Println("SPECIFICATION => ", config.ConfigTx.DbDataSourceName)
	if err != nil {
		log.Fatal(err.Error())
	}

	tx, err := db.SetupTestDb("../.env")
	if err != nil {
		log.Fatal(err)
	}
	defer tx.DB.Close()

	db.DefaultGameRepo = tx

	service.DefaultGameService = service.NewGameService(db.DefaultGameRepo)

	FiberEngine = fiber.New()

	FiberEngine.Use(logger.New(logger.Config{
		Format:     " ${pid} Time:${time} Status: ${status} - ${method} ${path}\n",
		TimeFormat: "02-Jan-2006 15:04:05",
		TimeZone:   "GMT+1",
	}))
	os.Exit(m.Run())
}

func deleteTestUsers(users []string) {

	for _, user := range users {
		for i := 1; i <= 3; i++ {

			db.DefaultGameRepo.DeleteTransactionByUsername(context.Background(), user)

		}
		db.DefaultGameRepo.DeleteWallet(context.Background(), user)

		db.DefaultGameRepo.DeleteUser(context.Background(), user)

	}
}

func createuserandcredit(user string, credit int) {
	dbuser, err := service.DefaultGameService.CreateUser(db.RegisterUserDto{
		Username:  user,
		Firstname: user + "-first",
		Lastname:  user + "-last",

		Password:        "test",
		ConfirmPassword: "test",
	})
	if err != nil {
		log.Fatal(err)
	}

	service.DefaultGameService.CreditWalletForWin(dbuser.Username, credit)
}

func createuser(user string) {
	_, err := service.DefaultGameService.CreateUser(db.RegisterUserDto{
		Username:  user,
		Firstname: user + "-first",
		Lastname:  user + "-last",

		Password:        "test",
		ConfirmPassword: "test",
	})
	if err != nil {
		log.Fatal(err)
	}

}
