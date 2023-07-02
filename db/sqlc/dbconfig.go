package db

import (
	"database/sql"
	"fmt"
)

func StartDb(DbDriverName string, DbSourceName string) {
	dbx, err := sql.Open(DbDriverName, DbSourceName)

	if err != nil {
		panic(fmt.Errorf("%s : %v", "cant connect to db", err))

	}

	DefaultStore = NewStore(dbx)

}
