package db

import (
	"context"
	"database/sql"
	"fmt"
)

const (
	CREDIT       = "CREDIT"
	DEBIT        = "DEBIT"
	UNSUCCESSFUL = "UNSUCCESSFUL"
)

type PGStore interface {
	Querier
	DebitWallet(ctx context.Context, arg UpdateWalletParams) error
	CreditWallet(ctx context.Context, arg UpdateWalletParams, win bool) error
	CreateUserTX(ctx context.Context, arg CreateUserParams) (User, Wallet, error)
}
type PGXDB struct {
	*Queries
	DB *sql.DB
}

func newPGXDB(db *sql.DB) *PGXDB {
	return &PGXDB{
		DB:      db,
		Queries: New(db),
	}
}

func NewPGXDB(drivername, sourcename string) (*PGXDB, error) {
	db, err := sql.Open(drivername, sourcename)
	if err != nil {
		return nil, fmt.Errorf("cant open database :  %+v", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("cant ping database :  %+v", err)
	}

	return newPGXDB(db), nil
}

func (s *PGXDB) execTx(ctx context.Context, fn func(*Queries) error) error {
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

func (s *PGXDB) DebitWallet(ctx context.Context, arg UpdateWalletParams) error {

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

func (s *PGXDB) CreditWallet(ctx context.Context, arg UpdateWalletParams, win bool) error {

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

func (s *PGXDB) CreateUserTX(ctx context.Context, arg CreateUserParams) (User, Wallet, error) {
	var err error
	var user User
	var wal Wallet
	err1 := s.execTx(ctx, func(q *Queries) error {

		user, err = q.CreateUser(ctx, arg)
		if err != nil {
			return fmt.Errorf("error creating user : %+v", err)
		}

		wal, err = q.CreateWallet(ctx, CreateWalletParams{
			UserID:   user.ID,
			Username: user.Username,
		})
		if err != nil {
			return fmt.Errorf("error creating wallet : %+v", err)
		}

		return nil
	})
	return user, wal, err1
}
