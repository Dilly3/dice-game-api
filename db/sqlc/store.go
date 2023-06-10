package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	execTx(ctx context.Context, fn func(*Queries) error) error
	DebitWallet(ctx context.Context, arg UpdateWalletParams) error
	CreditWallet(ctx context.Context, arg UpdateWalletParams, win bool) error
}

type PGXStore struct {
	*Queries
	DB *sql.DB
}

const (
	CREDIT       = "CREDIT"
	DEBIT        = "DEBIT"
	UNSUCCESSFUL = "UNSUCCESSFUL"
)

func NewStore(db *sql.DB) Store {
	return &PGXStore{
		DB:      db,
		Queries: New(db),
	}
}

func (s *PGXStore) execTx(ctx context.Context, fn func(*Queries) error) error {
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
	Amount   int32
	Username string
}

func (s *PGXStore) DebitWallet(ctx context.Context, arg UpdateWalletParams) error {

	err1 := s.execTx(ctx, func(q *Queries) error {
		var err error
		var wal Wallet

		wal, err = q.GetWalletByUsernameForUpdate(ctx, arg.Username)
		if err != nil {
			return err
		}
		if wal.Balance < arg.Balance {
			return fmt.Errorf("insufficient funds : %v", err)
		}

		err = q.UpdateWallet(ctx, UpdateWalletParams{
			Username: arg.Username,
			Balance:  wal.Balance - arg.Balance,
		})
		if err != nil {
			return err
		}

		_, err = q.CreateTransaction(ctx, CreateTransactionParams{
			UserID:          wal.UserID,
			Amount:          arg.Balance,
			TransactionType: DEBIT,
			Balance:         wal.Balance - arg.Balance,
			Username:        wal.Username,
		})

		if err != nil {
			return err
		}

		return err
	})

	return err1
}

func (s *PGXStore) CreditWallet(ctx context.Context, arg UpdateWalletParams, win bool) error {

	err1 := s.execTx(ctx, func(q *Queries) error {
		var err error
		var wal Wallet

		wal, err = q.GetWalletByUsernameForUpdate(ctx, arg.Username)
		if err != nil {
			return err
		}
		if !win {
			if arg.Balance != 155 {
				return fmt.Errorf("can only credit 155 sats : %v", err)
			}
			if wal.Balance >= 35 {
				return fmt.Errorf("you still have up to 35 sats : %v", err)
			}
		}

		err = q.UpdateWallet(ctx, UpdateWalletParams{
			Username: arg.Username,
			Balance:  wal.Balance + arg.Balance,
		})
		if err != nil {
			return err
		}

		_, err = q.CreateTransaction(ctx, CreateTransactionParams{
			UserID:          wal.UserID,
			Amount:          arg.Balance,
			TransactionType: CREDIT,
			Balance:         arg.Balance + wal.Balance,
			Username:        wal.Username,
		})

		if err != nil {
			return err
		}

		return err
	})

	return err1
}
