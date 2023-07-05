package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/dilly3/dice-game-api/models"
	"github.com/dilly3/dice-game-api/repository"
	"golang.org/x/crypto/bcrypt"
)

var once sync.Once
var DefaultGameService GameService

type GameService struct {
	Database repository.GameRepo
}

func newGameService(db repository.GameRepo) GameService {

	return GameService{
		Database: db,
	}
}

func GetGameService(db repository.GameRepo) GameService {
	once.Do(func() {
		DefaultGameService = newGameService(db)
	})
	return DefaultGameService

}
func (s *GameService) CreateUser(userData models.CreateUserParams) (models.User, error) {

	hashpassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, fmt.Errorf("cant hash password : %v", err)
	}
	usr := models.CreateUserParams{
		Firstname: userData.Firstname,
		Lastname:  userData.Lastname,
		Username:  userData.Username,
		Password:  string(hashpassword),
	}

	user, err := s.Database.CreateUser(context.Background(), usr)
	if err != nil {
		return models.User{}, fmt.Errorf("database error : %v", err)
	}

	s.Database.CreateWallet(context.Background(), models.CreateWalletParams{
		UserID:   user.ID,
		Username: user.Username,
	})
	return user, nil

}

func (s *GameService) GetAllUsers(ctx context.Context, arg models.ListUsersParams) ([]models.User, error) {
	return s.Database.ListUsers(ctx, arg)
}
func (s *GameService) GetUserByUsername(ctx context.Context, username string) (models.User, error) {
	return s.Database.GetUserByUsername(ctx, username)
}

func (s *GameService) GetWalletBalance(username string) (int32, string, error) {
	wallet, err := s.Database.GetWalletByUsername(context.Background(), username)
	if err != nil {
		return 0, "", fmt.Errorf("cant get wallet : %v", err)
	}
	return wallet.Balance, wallet.Assets, nil
}

func (s GameService) CreditWallet(username string, amount int32) error {
	err := s.Database.CreditWallet(context.Background(), models.UpdateWalletParams{
		Balance:  amount,
		Username: username,
	}, false)

	if err != nil {
		return err
	}
	return nil
}

func (s GameService) CreditWalletForWin(username string, amount int32) error {
	err := s.Database.CreditWallet(context.Background(), models.UpdateWalletParams{
		Balance:  amount,
		Username: username,
	}, true)

	if err != nil {
		return err
	}
	return nil
}

func (s GameService) DebitWallet(username string, amount int32) error {
	err := s.Database.DebitWallet(context.Background(), models.UpdateWalletParams{
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
	err = s.Database.UpdateUserGameMode(context.Background(), models.UpdateUserGameModeParams{
		Username: user.Username,
		GameMode: mode,
	})

	if err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func (s GameService) GetTransactionHistory(username string) ([]models.Transaction, error) {
	return s.Database.GetTransactionsByUsername(context.Background(), username)
}

func (s GameService) CreateTransaction(args models.CreateTransactionParams) (models.Transaction, error) {
	return s.Database.CreateTransaction(context.Background(), args)
}
