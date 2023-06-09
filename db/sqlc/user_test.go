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
		Email:     util.GenerateRandomEmail(6),
		Password:  util.GenerateRandomString(10),
	}
	user, err := StoreIntx.CreateUser(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, params.Firstname, user.Firstname)
	user, err = StoreIntx.GetUserByUsername(context.Background(), user.Username)
	wallet, err := StoreIntx.CreateWallet(context.Background(), CreateWalletParams{
		UserID:   user.ID,
		Username: user.Username,
	})
	if err != nil {
		panic(err)
	}
	return user, wallet
}

func TestCreateUser(t *testing.T) {
	CreateUser(t)
	//testQueries.DeleteAccount(context.Background(), user.Username)

}
