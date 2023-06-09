package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	Querier
	DB *sql.DB
}

const (
	CREDIT       = "CREDIT"
	DEBIT        = "DEBIT"
	UNSUCCESSFUL = "UNSUCCESSFUL"
)

func NewStore(db *sql.DB) *Store {
	return &Store{
		DB:      db,
		Querier: New(db),
	}
}

func (s *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			fmt.Println("Rollback error: ", rbErr)
		}
		return err
	}
	return tx.Commit()
}

type DebitWalletParam struct {
	Amount   int64
	Username string
}

func (s *Store) DebitWallet(ctx context.Context, arg DebitWalletParam) error {
	var err1 error
	err1 = s.execTx(ctx, func(q *Queries) error {
		var err error
		var wal Wallet

		wal, err = q.GetWalletByUsernameForUpdate(ctx, arg.Username)
		if err != nil {
			return err
		}
		if wal.Balance < arg.Amount {
			return fmt.Errorf("insufficient funds : %v", err)
		}

		err = q.UpdateWallet(ctx, UpdateWalletParams{
			Username: arg.Username,
			Balance:  wal.Balance - arg.Amount,
		})
		if err != nil {
			return err
		}

		_, err = q.CreateTransaction(ctx, CreateTransactionParams{
			UserID:          wal.UserID,
			Amount:          arg.Amount,
			TransactionType: DEBIT,
		})

		if err != nil {
			return err
		}

		return err
	})

	return err1
}

func (s *Store) CreditWallet(ctx context.Context, arg DebitWalletParam) error {
	var err1 error
	err1 = s.execTx(ctx, func(q *Queries) error {
		var err error
		var wal Wallet

		wal, err = q.GetWalletByUsernameForUpdate(ctx, arg.Username)
		if err != nil {
			return err
		}
		if arg.Amount > 155 {
			return fmt.Errorf("cant credit above 155sats : %v", err)
		}
		if wal.Balance >= 35 {
			return fmt.Errorf("still have enough money : %v", err)
		}

		err = q.UpdateWallet(ctx, UpdateWalletParams{
			Username: arg.Username,
			Balance:  wal.Balance + arg.Amount,
		})
		if err != nil {
			return err
		}

		_, err = q.CreateTransaction(ctx, CreateTransactionParams{
			UserID:          wal.UserID,
			Amount:          arg.Amount,
			TransactionType: CREDIT,
		})

		if err != nil {
			return err
		}

		return err
	})

	return err1
}
