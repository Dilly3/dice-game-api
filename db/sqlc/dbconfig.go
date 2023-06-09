package db

import (
	"database/sql"
	"fmt"
)

func StartDb(DbDriverName string, DbSourceName string) Store {
	dbx, err := sql.Open(DbDriverName, DbSourceName)

	if err != nil {
		panic(fmt.Errorf("%s : %v", "cant connect to db", err))

	}

	StoreIntx := NewStore(dbx)
	return StoreIntx
}
