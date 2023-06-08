package config

import (
	"database/sql"
)

type Configuration struct {
	DbDriverName         string
	DbDataSourceNameTest string
	DbDataSourceName     string
	Db                   *sql.DB
	IsGameInSession      bool
}

var ConfigTx Configuration
