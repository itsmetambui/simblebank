package db

import (
	"context"
	"testing"
	"time"

	"github.com/ddosify/go-faker/faker"
	"github.com/stretchr/testify/assert"
)

func createRandomTransfer(t *testing.T, fromAccount Account, toAccount Account) Transfer {
	faker := faker.NewFaker()
	arg := CreateTransferParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Amount:        int64(faker.RandomInt()),
	}
	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, transfer)

	assert.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	assert.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	assert.Equal(t, arg.Amount, transfer.Amount)

	assert.NotZero(t, transfer.ID)
	assert.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	createRandomTransfer(t, fromAccount, toAccount)
}

func TestGetTransfer(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	transfer1 := createRandomTransfer(t, fromAccount, toAccount)
	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, transfer2)

	assert.Equal(t, transfer1.ID, transfer2.ID)
	assert.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	assert.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	assert.Equal(t, transfer1.Amount, transfer2.Amount)
	assert.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)
}

func TestListTransfers(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		createRandomTransfer(t, fromAccount, toAccount)
	}

	arg := ListTransfersParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Limit:         5,
		Offset:        5,
	}
	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	assert.NoError(t, err)
	assert.Len(t, transfers, 5)

	for _, transfer := range transfers {
		assert.NotEmpty(t, transfer)
		assert.Equal(t, arg.FromAccountID, transfer.FromAccountID)
		assert.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	}
}
