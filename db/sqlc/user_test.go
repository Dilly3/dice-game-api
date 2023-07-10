package db

import (
	"context"
	"testing"

	"github.com/dilly3/dice-game-api/util"
	"github.com/stretchr/testify/require"
)

func CreateUser(t *testing.T) (User, Wallet) {
	params := CreateUserParams{
		Firstname: util.GenerateRandomString(8),
		Lastname:  util.GenerateRandomString(9),
		Username:  util.GenerateRandomUsername(5),
		Password:  util.GenerateRandomString(10),
	}
	user, wal, err := DefaultGameRepo.CreateUserTX(context.Background(), params)

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.NotEmpty(t, wal)
	require.Equal(t, params.Username, user.Username)
	require.Equal(t, params.Username, wal.Username)
	require.Equal(t, params.Firstname, user.Firstname)
	require.Equal(t, wal.Balance, 0)
	require.Equal(t, params.Firstname, user.Firstname)
	user, err = DefaultGameRepo.GetUserByUsername(context.Background(), user.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	return user, wal
}

func TestCreateUser(t *testing.T) {
	user, wal := CreateUser(t)
	defer func() error {
		err := DefaultGameRepo.DeleteWallet(context.Background(), wal.Username)
		if err != nil {
			return err
		}

		err = DefaultGameRepo.DeleteUser(context.Background(), user.Username)
		if err != nil {
			return err
		}
		return nil
	}()
	//testQueries.DeleteAccount(context.Background(), user.Username)

}
