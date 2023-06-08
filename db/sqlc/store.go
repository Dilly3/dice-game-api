package db

import "database/sql"

type Store struct {
	Querier
	DB *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		DB:      db,
		Querier: New(db),
	}
}
