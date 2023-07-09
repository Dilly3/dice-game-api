package db

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
	Amount int `json:"amount"`
}
type UpdateWalletDto struct {
	Username string `json:"username"`
	Amount   int    `json:"amount"`
}
