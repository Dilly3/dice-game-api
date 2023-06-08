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
	return user, nil

}

func (s *UserService) GetAllUsers(ctx context.Context, arg db.ListUsersParams) ([]db.User, error) {
	return s.Database.ListUsers(ctx, arg)
}
