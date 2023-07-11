package db

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type RegisterUserDto struct {
	Firstname       string `json:"firstname" binding:"required"`
	Lastname        string `json:"lastname" binding:"required"`
	Username        string `json:"username" binding:"required"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

type LoginDto struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreateWalletDto struct {
	Amount int `json:"amount"`
}
type UpdateWalletDto struct {
	Username string `json:"username"`
	Amount   int    `json:"amount"`
}

type IUser interface {
	HashPassword(user *User)
	CompareHashAndPassword(user User, password string) error
	VerifyUserData(user *RegisterUserDto) error
}

func HashPassword(user *User) {
	hashpassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	user.Password = string(hashpassword)
}

func CompareHashAndPassword(user User, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}

func VerifyUserData(user *RegisterUserDto) error {
	if user.Firstname == "" || user.Lastname == "" || user.Username == "" || user.Password == "" {
		return errors.New("some fields missing in user data")
	}
	return nil
}
