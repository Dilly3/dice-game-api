package repository

import (
	"context"

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
