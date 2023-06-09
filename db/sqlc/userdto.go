package db

import "golang.org/x/crypto/bcrypt"

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

func (user *User) HashPassword() {
	hashpassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	user.Password = string(hashpassword)
}

func (user *User) CompareHashAndPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}
