package service

import (
	"context"
	"fmt"

	db "github.com/dilly3/dice-game-api/db/sqlc"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Database *db.Store
}

func NewUserService(db *db.Store) *UserService {
	return &UserService{
		Database: db,
	}
}

func (s *UserService) CreateUser(userData db.CreateUserParams) (db.User, error) {

	hashpassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if err != nil {
		return db.User{}, fmt.Errorf("cant hash password : %v", err)
	}
	usr := db.CreateUserParams{
		Firstname: userData.Firstname,
		Lastname:  userData.Lastname,
		Username:  userData.Username,
		Email:     userData.Email,
		Password:  string(hashpassword),
	}

	user, err := s.Database.CreateUser(context.Background(), usr)
	if err != nil {
		return db.User{}, fmt.Errorf("database error : %v", err)
	}

	s.Database.CreateWallet(context.Background(), db.CreateWalletParams{
		UserID:   user.ID,
		Username: user.Username,
	})
	return user, nil

}

func (s *UserService) GetAllUsers(ctx context.Context, arg db.ListUsersParams) ([]db.User, error) {
	return s.Database.ListUsers(ctx, arg)
}

func (s *UserService) GetWalletBalance(username string) (int64, error) {
	wallet, err := s.Database.Querier.GetWalletByUsername(context.Background(), username)
	if err != nil {
		return 0, fmt.Errorf("canr get wallet : %v", err)
	}
	return wallet.Balance, nil
}

func (s UserService) CreditWallet(username string, amount int64) error {
	err := s.Database.CreditWallet(context.Background(), db.DebitWalletParam{
		Amount:   amount,
		Username: username,
	})

	if err != nil {
		return err
	}
	return nil
}

func (s UserService) DebitWallet(username string, amount int64) error {
	err := s.Database.DebitWallet(context.Background(), db.DebitWalletParam{
		Amount:   amount,
		Username: username,
	})

	if err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}
