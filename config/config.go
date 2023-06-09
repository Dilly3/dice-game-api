package config

import (
	"database/sql"
	"fmt"
	"math/rand"
)

type Configuration struct {
	DbDriverName         string
	DbDataSourceNameTest string
	DbDataSourceName     string
	Db                   *sql.DB
	IsGameInSession      bool
	RollNumber1          int32
	RollNumber2          int32
	LuckyNumber          int32
	NumberOfTrials       int32
}

var ConfigTx Configuration

type RollResult struct {
	RollNumber1   int32 `json:"rollNumber1"`
	RollNumber2   int32 `json:"rollNumber2"`
	LuckyNumber   int32 `json:"luckyNumber"`
	NumberOfTries int32 `json:"numberOfTries"`
}

func StartGame() {

	randNum := rand.Int31n(11) + 2
	ConfigTx.IsGameInSession = true
	ConfigTx.LuckyNumber = randNum
	ConfigTx.NumberOfTrials = 10

}

func RollDice() (*RollResult, error) {
	if ConfigTx.NumberOfTrials == 0 {
		ConfigTx.IsGameInSession = false

		return nil, fmt.Errorf("%s", "game over")
	}
	rollResult := &RollResult{}
	if !ConfigTx.IsGameInSession {
		return rollResult, fmt.Errorf("%s", "there is no game in session")
	}

	num := rand.Int31n(6) + 1
	ConfigTx.RollNumber1 = num
	rollResult.RollNumber1 = num
	num2 := rand.Int31n(6) + 1
	ConfigTx.RollNumber2 = num2
	rollResult.RollNumber2 = num2
	rollResult.LuckyNumber = ConfigTx.LuckyNumber
	ConfigTx.NumberOfTrials = ConfigTx.NumberOfTrials - 1
	rollResult.NumberOfTries = ConfigTx.NumberOfTrials

	return rollResult, nil

}
