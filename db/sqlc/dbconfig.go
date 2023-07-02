package db

import (
	"database/sql"
	"fmt"
)

var DefaultStore Store

var StartDb = func(DbDriverName string, DbSourceName string) Store {
	db := opendb(DbDriverName, DbSourceName)
	store := NewStore(db)
	return setDefaultStore(store)

}

var opendb = func(DbDriverName string, DbSourceName string) *sql.DB {
	dbx, err := sql.Open(DbDriverName, DbSourceName)

	if err != nil {
		panic(fmt.Errorf("%s : %v", "cant connect to db", err))

	}

	return dbx

}

// set default store
func setDefaultStore(db Store) Store {

	DefaultStore = db
	return DefaultStore

}
