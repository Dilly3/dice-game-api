package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetBalance(t *testing.T) {
	user, wallet := CreateUser(t)
	defer func() error {
		err := DefaultGameRepo.DeleteWallet(context.Background(), wallet.Username)
		if err != nil {
			return err
		}

		err = DefaultGameRepo.DeleteUser(context.Background(), user.Username)
		if err != nil {
			return err
		}
		return nil
	}()
	wallet, err := DefaultGameRepo.GetWalletByUsername(context.Background(), user.Username)
	require.NoError(t, err)
	require.Equal(t, wallet.Balance, 0)

	DefaultGameRepo.UpdateWallet(context.Background(), UpdateWalletParams{
		Username: user.Username,
		Balance:  100,
	})
	wallet, err = DefaultGameRepo.GetWalletByUsername(context.Background(), user.Username)
	require.NoError(t, err)
	require.Equal(t, wallet.Balance, 100)
	require.Equal(t, wallet.Username, user.Username)
	require.Equal(t, wallet.UserID, user.ID)

}
