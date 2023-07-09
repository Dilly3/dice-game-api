package service

import (
	"context"
	"fmt"
	"sync"

	db "github.com/dilly3/dice-game-api/db/sqlc"
	"golang.org/x/crypto/bcrypt"
)

type IGameService interface {
	CreateUser(userData db.CreateUserParams) (db.User, error)
	GetUserByUsername(ctx context.Context, username string) (db.User, error)
	GetAllUsers(ctx context.Context, arg db.ListUsersParams) ([]db.User, error)
	CreditWalletForWin(username string, amount int) error
	CreditWallet(username string, amount int) error
	DebitWallet(username string, amount int) error
	UpdateGameMode(username string, mode bool) error
	GetTransactionHistory(username string) ([]db.Transaction, error)
	CreateTransaction(args db.CreateTransactionParams) (db.Transaction, error)
	GetWalletBalance(username string) (int, string, error)
}

var DefaultGameService IGameService

var once sync.Once

type GameService struct {
	Database db.IGameRepo
}

func NewGameService(db db.IGameRepo) *GameService {

	return &GameService{
		Database: db,
	}
}

func GetGameService(db db.IGameRepo) IGameService {
	once.Do(func() {
		DefaultGameService = NewGameService(db)
	})
	return DefaultGameService

}
func (s *GameService) CreateUser(userData db.CreateUserParams) (db.User, error) {

	hashpassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if err != nil {
		return db.User{}, fmt.Errorf("cant hash password : %v", err)
	}
	usr := db.CreateUserParams{
		Firstname: userData.Firstname,
		Lastname:  userData.Lastname,
		Username:  userData.Username,
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

func (s *GameService) GetAllUsers(ctx context.Context, arg db.ListUsersParams) ([]db.User, error) {
	return s.Database.ListUsers(ctx, arg)
}
func (s *GameService) GetUserByUsername(ctx context.Context, username string) (db.User, error) {
	return s.Database.GetUserByUsername(ctx, username)
}

func (s *GameService) GetWalletBalance(username string) (int, string, error) {
	wallet, err := s.Database.GetWalletByUsername(context.Background(), username)
	if err != nil {
		return 0, "", fmt.Errorf("cant get wallet : %v", err)
	}
	return wallet.Balance, wallet.Assets, nil
}

func (s GameService) CreditWallet(username string, amount int) error {
	err := s.Database.CreditWallet(context.Background(), db.UpdateWalletParams{
		Balance:  amount,
		Username: username,
	}, false)

	if err != nil {
		return err
	}
	return nil
}

func (s GameService) CreditWalletForWin(username string, amount int) error {
	err := s.Database.CreditWallet(context.Background(), db.UpdateWalletParams{
		Balance:  amount,
		Username: username,
	}, true)

	if err != nil {
		return err
	}
	return nil
}

func (s GameService) DebitWallet(username string, amount int) error {
	err := s.Database.DebitWallet(context.Background(), db.UpdateWalletParams{
		Balance:  amount,
		Username: username,
	})

	if err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func (s GameService) UpdateGameMode(username string, mode bool) error {

	user, err := s.Database.GetUserForUpdate(context.Background(), username)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	err = s.Database.UpdateUserGameMode(context.Background(), db.UpdateUserGameModeParams{
		Username: user.Username,
		GameMode: mode,
	})

	if err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func (s GameService) GetTransactionHistory(username string) ([]db.Transaction, error) {
	return s.Database.GetTransactionsByUsername(context.Background(), username)
}

func (s GameService) CreateTransaction(args db.CreateTransactionParams) (db.Transaction, error) {
	return s.Database.CreateTransaction(context.Background(), args)
}
