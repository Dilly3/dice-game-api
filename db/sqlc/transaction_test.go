package db

import (
	"context"

	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetTransactionByUsername(t *testing.T) {
	user, wallet := CreateUserSample(t)

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

	DefaultGameRepo.CreditWallet(context.Background(), UpdateWalletParams{
		Username: user.Username,
		Balance:  100,
	}, true)

	DefaultGameRepo.DebitWallet(context.Background(), UpdateWalletParams{
		Username: user.Username,
		Balance:  10,
	})

	tranx, err := DefaultGameRepo.GetTransactionsByUsername(context.Background(), GetTransactionsByUsernameParams{Username: user.Username, Limit: 3})
	if err != nil {
		t.Fail()
	}

	require.NotNil(t, tranx)
	require.Equal(t, tranx[0].Amount, 10)
	require.Equal(t, tranx[0].Balance, 90)
	require.Equal(t, tranx[0].Username, user.Username)
	require.Equal(t, tranx[1].TransactionType, "CREDIT")
	require.Equal(t, tranx[0].TransactionType, "DEBIT")
	require.Equal(t, tranx[1].Balance, 100)

}
