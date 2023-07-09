package db

import "context"

type IGameRepo interface {
	CreateTransaction(ctx context.Context, arg CreateTransactionParams) (Transaction, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	CreateWallet(ctx context.Context, arg CreateWalletParams) (Wallet, error)
	DeleteUser(ctx context.Context, username string) error
	GetTransaction(ctx context.Context, arg GetTransactionParams) (Transaction, error)
	GetTransactionsByUsername(ctx context.Context, username string) ([]Transaction, error)
	GetUser(ctx context.Context, username string) (User, error)
	GetUserByUsername(ctx context.Context, username string) (User, error)
	GetUserForUpdate(ctx context.Context, username string) (User, error)
	GetWalletByUsername(ctx context.Context, username string) (Wallet, error)
	GetWalletByUsernameForUpdate(ctx context.Context, username string) (Wallet, error)
	ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error)
	UpdateTransaction(ctx context.Context, arg UpdateTransactionParams) error
	UpdateUserGameMode(ctx context.Context, arg UpdateUserGameModeParams) error
	UpdateWallet(ctx context.Context, arg UpdateWalletParams) error
	DebitWallet(ctx context.Context, arg UpdateWalletParams) error
	CreditWallet(ctx context.Context, arg UpdateWalletParams, win bool) error
}

var DefaultGameRepo IGameRepo
