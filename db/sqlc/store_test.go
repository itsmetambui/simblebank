package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// run n concurrent tranfer transactions
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		assert.NoError(t, err)

		result := <-results
		assert.NotEmpty(t, result)

		// Check the transfer details
		assert.NotEmpty(t, result.Transfer)
		assert.Equal(t, account1.ID, result.Transfer.FromAccountID)
		assert.Equal(t, account2.ID, result.Transfer.ToAccountID)
		assert.Equal(t, amount, result.Transfer.Amount)
		assert.NotZero(t, result.Transfer.ID)
		assert.NotZero(t, result.Transfer.CreatedAt)

		// Check the entries
		assert.NotEmpty(t, result.FromEntry)
		assert.Equal(t, account1.ID, result.FromEntry.AccountID)
		assert.Equal(t, -amount, result.FromEntry.Amount)
		assert.NotZero(t, result.FromEntry.ID)
		assert.NotZero(t, result.FromEntry.CreatedAt)

		assert.NotEmpty(t, result.ToEntry)
		assert.Equal(t, account2.ID, result.ToEntry.AccountID)
		assert.Equal(t, amount, result.ToEntry.Amount)
		assert.NotZero(t, result.ToEntry.ID)
		assert.NotZero(t, result.ToEntry.CreatedAt)

		// Check the accounts
		assert.NotEmpty(t, result.FromAccount)
		assert.Equal(t, account1.ID, result.FromAccount.ID)

		assert.NotEmpty(t, result.ToAccount)
		assert.Equal(t, account2.ID, result.ToAccount.ID)

		diff1 := account1.Balance - result.FromAccount.Balance
		diff2 := result.ToAccount.Balance - account2.Balance
		assert.Equal(t, diff1, diff2)
		assert.True(t, diff1 > 0)
		assert.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		assert.True(t, k >= 1 && k <= n)
		assert.NotContains(t, existed, k)
		existed[k] = true
	}

	// Check the final updated balances
	updatedFromAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	assert.NoError(t, err)

	updatedFromAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	assert.NoError(t, err)

	assert.Equal(t, account1.Balance-int64(n)*amount, updatedFromAccount1.Balance)
	assert.Equal(t, account2.Balance+int64(n)*amount, updatedFromAccount2.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// run n concurrent tranfer transactions
	n := 10
	amount := int64(10)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		// Reverse the account IDs in odd transfers to create a deadlock
		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})

			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		assert.NoError(t, err)
	}

	// Check the final updated balances
	updatedFromAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	assert.NoError(t, err)

	updatedFromAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	assert.NoError(t, err)

	assert.Equal(t, account1.Balance, updatedFromAccount1.Balance)
	assert.Equal(t, account2.Balance, updatedFromAccount2.Balance)
}
