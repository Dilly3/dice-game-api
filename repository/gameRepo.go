package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dilly3/dice-game-api/models"
)

type GameRepo interface {
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

// repo instance
var DefaultRepoInstance GameRepo

var StartDb = func(DbDriverName string, DbSourceName string, initdb func(*sql.DB) GameRepo) GameRepo {
	dbx := opendb(DbDriverName, DbSourceName)
	gamerepo := initdb(dbx)
	//set repo instance
	return gamerepo

}

var opendb = func(DbDriverName string, DbSourceName string) *sql.DB {
	dbx, err := sql.Open(DbDriverName, DbSourceName)

	if err != nil {
		panic(fmt.Errorf("%s : %v", "cant connect to db", err))

	}

	return dbx

}

// set default repo instance
func setRepoInstance(db GameRepo) GameRepo {

	DefaultRepoInstance = db
	return DefaultRepoInstance

}
