// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0

package db

import (
	"context"
)

type Querier interface {
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
	UpdateUserGameMode(ctx context.Context, arg UpdateUserGameModeParams) error
	UpdateWallet(ctx context.Context, arg UpdateWalletParams) error
}

var _ Querier = (*Queries)(nil)
