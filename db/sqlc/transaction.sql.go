// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: transaction.sql

package db

import (
	"context"
)

const createTransaction = `-- name: CreateTransaction :one
INSERT INTO transactions (
    user_id, amount, balance , transaction_type, username
) VALUES (
  $1, $2 , $3 , $4 , $5
)
RETURNING id, user_id, username, amount, balance, transaction_type, created_at
`

type CreateTransactionParams struct {
	UserID          int64  `json:"user_id"`
	Amount          int  `json:"amount"`
	Balance         int  `json:"balance"`
	TransactionType string `json:"transaction_type"`
	Username        string `json:"username"`
}

func (q *Queries) CreateTransaction(ctx context.Context, arg CreateTransactionParams) (Transaction, error) {
	row := q.db.QueryRowContext(ctx, createTransaction,
		arg.UserID,
		arg.Amount,
		arg.Balance,
		arg.TransactionType,
		arg.Username,
	)
	var i Transaction
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Username,
		&i.Amount,
		&i.Balance,
		&i.TransactionType,
		&i.CreatedAt,
	)
	return i, err
}

const getTransaction = `-- name: GetTransaction :one
SELECT id, user_id, username, amount, balance, transaction_type, created_at FROM transactions
WHERE user_id = $1
AND transaction_type = $2
`

type GetTransactionParams struct {
	UserID          int64  `json:"user_id"`
	TransactionType string `json:"transaction_type"`
}

func (q *Queries) GetTransaction(ctx context.Context, arg GetTransactionParams) (Transaction, error) {
	row := q.db.QueryRowContext(ctx, getTransaction, arg.UserID, arg.TransactionType)
	var i Transaction
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Username,
		&i.Amount,
		&i.Balance,
		&i.TransactionType,
		&i.CreatedAt,
	)
	return i, err
}

const getTransactionsByUsername = `-- name: GetTransactionsByUsername :many
SELECT id, user_id, username, amount, balance, transaction_type, created_at FROM transactions
WHERE username = $1
ORDER BY created_at DESC
`

func (q *Queries) GetTransactionsByUsername(ctx context.Context, username string) ([]Transaction, error) {
	rows, err := q.db.QueryContext(ctx, getTransactionsByUsername, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Transaction
	for rows.Next() {
		var i Transaction
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Username,
			&i.Amount,
			&i.Balance,
			&i.TransactionType,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateTransaction = `-- name: UpdateTransaction :exec
UPDATE transactions
  set balance = $2 ,
  amount = $3
WHERE username = $1
`

type UpdateTransactionParams struct {
	Username string `json:"username"`
	Balance  int `json:"balance"`
	Amount   int  `json:"amount"`
}

func (q *Queries) UpdateTransaction(ctx context.Context, arg UpdateTransactionParams) error {
	_, err := q.db.ExecContext(ctx, updateTransaction, arg.Username, arg.Balance, arg.Amount)
	return err
}
