package service

import (
	"context"

	db "github.com/dilly3/dice-game-api/db/sqlc"
)

type TransactionService struct {
	Database db.Store
}

func NewTransactionService(db db.Store) *TransactionService {
	return &TransactionService{
		Database: db,
	}
}

func (s TransactionService) GetTransactionHistory(username string) ([]db.Transaction, error) {
	return s.Database.GetTransactionsByUsername(context.Background(), username)
}

func (s TransactionService) CreateTransaction(args db.CreateTransactionParams) (db.Transaction, error) {
	return s.Database.CreateTransaction(context.Background(), args)
}
