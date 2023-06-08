package service

import db "github.com/dilly3/dice-game-api/db/sqlc"

type TransactionService struct {
	Database *db.Store
}

func NewTransactionService(db *db.Store) *TransactionService {
	return &TransactionService{
		Database: db,
	}
}
