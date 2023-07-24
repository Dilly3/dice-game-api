package db

import (
	"log"

	"github.com/dilly3/dice-game-api/config"
)

type databaseFactory struct {
}

type Database string

const (
	POSTGRES Database = "postgres"
)

func NewDatabaseFactory() *databaseFactory {
	return &databaseFactory{}
}

func (dbf *databaseFactory) GetDatabaseInstance(database Database, config *config.Configuration) IGameRepo {
	switch database {
	case "postgres":
		pgx, err := NewPGXDB(config.DbDriverName, config.DbDataSourceName)
		if err != nil {
			log.Fatal("error from database factory: ", err.Error())
		}
		return pgx
	default:
		return nil
	}
}
