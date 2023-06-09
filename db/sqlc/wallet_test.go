package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetBalance(t *testing.T) {
	user, wallet := CreateUser(t)
	wallet, err := StoreIntx.GetWalletByUsername(context.Background(), user.Username)
	require.NoError(t, err)
	require.Equal(t, wallet.Balance, int32(0))

	StoreIntx.UpdateWallet(context.Background(), UpdateWalletParams{
		Username: user.Username,
		Balance:  100,
	})
	wallet, err = StoreIntx.GetWalletByUsername(context.Background(), user.Username)
	require.NoError(t, err)
	require.Equal(t, wallet.Balance, int32(100))
	require.Equal(t, wallet.Username, user.Username)
	require.Equal(t, wallet.UserID, user.ID)

}
