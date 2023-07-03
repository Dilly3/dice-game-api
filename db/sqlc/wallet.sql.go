// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: wallet.sql

package db

import (
	"context"

	"github.com/dilly3/dice-game-api/models"
)

const createWallet = `-- name: CreateWallet :one
INSERT INTO wallets (
  user_id, username
) VALUES (
  $1, $2 
)
RETURNING id, user_id, username, balance, assets, updated_at
`



func (q *Queries) CreateWallet(ctx context.Context, arg models.CreateWalletParams) (models.Wallet, error) {
	row := q.db.QueryRowContext(ctx, createWallet, arg.UserID, arg.Username)
	var i models.Wallet
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Username,
		&i.Balance,
		&i.Assets,
		&i.UpdatedAt,
	)
	return i, err
}

const getWalletByUsername = `-- name: GetWalletByUsername :one
SELECT id, user_id, username, balance, assets, updated_at FROM wallets
WHERE username = $1
`

func (q *Queries) GetWalletByUsername(ctx context.Context, username string) (models.Wallet, error) {
	row := q.db.QueryRowContext(ctx, getWalletByUsername, username)
	var i models.Wallet
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Username,
		&i.Balance,
		&i.Assets,
		&i.UpdatedAt,
	)
	return i, err
}

const getWalletByUsernameForUpdate = `-- name: GetWalletByUsernameForUpdate :one
SELECT id, user_id, username, balance, assets, updated_at FROM wallets
WHERE username = $1
FOR UPDATE
`

func (q *Queries) GetWalletByUsernameForUpdate(ctx context.Context, username string) (models.Wallet, error) {
	row := q.db.QueryRowContext(ctx, getWalletByUsernameForUpdate, username)
	var i models.Wallet
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Username,
		&i.Balance,
		&i.Assets,
		&i.UpdatedAt,
	)
	return i, err
}

const updateWallet = `-- name: UpdateWallet :exec
UPDATE wallets
  set balance = $2
WHERE username = $1
`



func (q *Queries) UpdateWallet(ctx context.Context, arg models.UpdateWalletParams) error {
	_, err := q.db.ExecContext(ctx, updateWallet, arg.Username, arg.Balance)
	return err
}
