package postgres

import (
	"context"
	"testing"
	"tgBank/models"
	"tgBank/utils"
	"time"

	"github.com/stretchr/testify/require"
)

func CreateRandomAccount(t *testing.T) models.Account {
	arg := models.CreateAccountParams{
		Owner:    utils.RandomOwner(),
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}
	acc, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, acc)

	require.Equal(t, arg.Owner, acc.Owner)
	require.Equal(t, arg.Balance, acc.Balance)
	require.Equal(t, arg.Currency, acc.Currency)
	return acc
}

func TestCreateRandomAccount(t *testing.T) {
	CreateRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	acc := CreateRandomAccount(t)

	acc2, err := testQueries.GetAccount(context.Background(), acc.ID)
	require.NoError(t, err)

	require.Equal(t, acc.ID, acc2.ID)
	require.Equal(t, acc.Owner, acc2.Owner)
	require.Equal(t, acc.Balance, acc2.Balance)
	require.WithinDuration(t, acc.CreatedAt, acc2.CreatedAt, 1*time.Microsecond)
}

func TestUpdateAccount(t *testing.T) {
	acc := CreateRandomAccount(t)

	var argUpd = models.UpdateAccountParams{
		ID:      acc.ID,
		Balance: acc.Balance,
	}

	acc2, err := testQueries.UpdateAccount(context.Background(), argUpd)
	require.NoError(t, err)
	require.NotEmpty(t, acc2)

	require.Equal(t, acc.ID, acc2.ID)
	require.Equal(t, acc.Owner, acc2.Owner)
	require.Equal(t, acc.Balance, argUpd.Balance, acc2.Balance)
	require.WithinDuration(t, acc.CreatedAt, acc2.CreatedAt, 1*time.Microsecond)

}

func TestListAccounts(t *testing.T) {
	var acc []models.Account
	for i := 0; i < 10; i++ {
		acc = append(acc, CreateRandomAccount(t))
	}
	arg := models.ListAccountsParams{
		Offset: int32(acc[0].ID - 1),
		Limit:  10,
	}
	acc2, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	for i := 0; i < len(acc); i++ {
		require.Equal(t, acc[i].ID, acc2[i].ID)
		require.Equal(t, acc[i].Owner, acc2[i].Owner)
		require.Equal(t, acc[i].Balance, acc2[i].Balance)
		require.WithinDuration(t, acc[i].CreatedAt, acc2[i].CreatedAt, 1*time.Microsecond)
	}
}
