package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dilly3/dice-game-api/models"
	"github.com/dilly3/dice-game-api/repository"
)

// type Store interface {
// 	Querier
// 	DebitWallet(ctx context.Context, arg UpdateWalletParams) error
// 	CreditWallet(ctx context.Context, arg UpdateWalletParams, win bool) error
// }

type PGXStore struct {
	*Queries
	DB *sql.DB
}

const (
	CREDIT       = "CREDIT"
	DEBIT        = "DEBIT"
	UNSUCCESSFUL = "UNSUCCESSFUL"
)

type Store interface {
	CreateTransaction(ctx context.Context, arg models.CreateTransactionParams) (models.Transaction, error)
	CreateUser(ctx context.Context, arg models.CreateUserParams) (models.User, error)
	CreateWallet(ctx context.Context, arg models.CreateWalletParams) (models.Wallet, error)
	DeleteUser(ctx context.Context, username string) error
	GetTransaction(ctx context.Context, arg models.GetTransactionParams) (models.Transaction, error)
	GetTransactionsByUsername(ctx context.Context, username string) ([]models.Transaction, error)
	GetUser(ctx context.Context, username string) (models.User, error)
	GetUserByUsername(ctx context.Context, username string) (models.User, error)
	GetUserForUpdate(ctx context.Context, username string) (models.User, error)
	GetWalletByUsername(ctx context.Context, username string) (models.Wallet, error)
	GetWalletByUsernameForUpdate(ctx context.Context, username string) (models.Wallet, error)
	ListUsers(ctx context.Context, arg models.ListUsersParams) ([]models.User, error)
	UpdateTransaction(ctx context.Context, arg models.UpdateTransactionParams) error
	UpdateUserGameMode(ctx context.Context, arg models.UpdateUserGameModeParams) error
	UpdateWallet(ctx context.Context, arg models.UpdateWalletParams) error
	DebitWallet(ctx context.Context, arg models.UpdateWalletParams) error
	CreditWallet(ctx context.Context, arg models.UpdateWalletParams, win bool) error
}

func NewStore(db *sql.DB) repository.GameRepo {
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

func (s *PGXStore) DebitWallet(ctx context.Context, arg models.UpdateWalletParams) error {

	err1 := s.execTx(ctx, func(q *Queries) error {
		var err error
		var wal models.Wallet

		wal, err = q.GetWalletByUsernameForUpdate(ctx, arg.Username)
		if err != nil {
			return err
		}
		if wal.Balance < arg.Balance {
			return fmt.Errorf("insufficient funds : %v", err)
		}

		err = q.UpdateWallet(ctx, models.UpdateWalletParams{
			Username: arg.Username,
			Balance:  wal.Balance - arg.Balance,
		})
		if err != nil {
			return err
		}

		_, err = q.CreateTransaction(ctx, models.CreateTransactionParams{
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

func (s *PGXStore) CreditWallet(ctx context.Context, arg models.UpdateWalletParams, win bool) error {

	err1 := s.execTx(ctx, func(q *Queries) error {
		var err error
		var wal models.Wallet

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

		err = q.UpdateWallet(ctx, models.UpdateWalletParams{
			Username: arg.Username,
			Balance:  wal.Balance + arg.Balance,
		})
		if err != nil {
			return err
		}

		_, err = q.CreateTransaction(ctx, models.CreateTransactionParams{
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
