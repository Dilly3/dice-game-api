package models

type UpdateTransactionParams struct {
	Username string `json:"username"`
	Balance  int32  `json:"balance"`
	Amount   int32  `json:"amount"`
}

type GetTransactionParams struct {
	UserID          int64  `json:"user_id"`
	TransactionType string `json:"transaction_type"`
}

type CreateTransactionParams struct {
	UserID          int64  `json:"user_id"`
	Amount          int32  `json:"amount"`
	Balance         int32  `json:"balance"`
	TransactionType string `json:"transaction_type"`
	Username        string `json:"username"`
}

type DebitWalletParam struct {
	Amount   int32
	Username string
}

type CreateUserParams struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

type ListUsersParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type UpdateUserGameModeParams struct {
	Username string `json:"username"`
	GameMode bool   `json:"game_mode"`
}

type UpdateWalletParams struct {
	Username string `json:"username"`
	Balance  int32  `json:"balance"`
}

type CreateWalletParams struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
}

type RegisterUserDto struct {
	Firstname       string `json:"firstname"`
	Lastname        string `json:"lastname"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type LoginDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateWalletDto struct {
	Amount int32 `json:"amount"`
}
type UpdateWalletDto struct {
	Username string `json:"username"`
	Amount   int32  `json:"amount"`
}
