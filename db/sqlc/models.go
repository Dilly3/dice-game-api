// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0

package db

import (
	"time"
)

type Transaction struct {
	ID              int64     `json:"id"`
	UserID          int64     `json:"user_id"`
	Username        string    `json:"username"`
	Amount          int32     `json:"amount"`
	TransactionType string    `json:"transaction_type"`
	CreatedAt       time.Time `json:"created_at"`
}

type User struct {
	ID        int64     `json:"id"`
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	Username  string    `json:"username"`
	GameMode  bool      `json:"game_mode"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type Wallet struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Username  string    `json:"username"`
	Balance   int32     `json:"balance"`
	Assets    string    `json:"assets"`
	UpdatedAt time.Time `json:"updated_at"`
}
