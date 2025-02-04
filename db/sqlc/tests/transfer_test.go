package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgtype"
	db "github.com/riad/banksystemendtoend/db/sqlc"
	"github.com/riad/banksystemendtoend/db/sqlc/transaction"
	setup "github.com/riad/banksystemendtoend/util/db"
	"github.com/riad/banksystemendtoend/util/schemas"
	"github.com/stretchr/testify/require"
)

func TestTransfer(t *testing.T) {
	var err error

	sender := createRandomAccount(t)
	receiver := createRandomAccount(t)

	fmt.Println(">> before:", sender.Balance, receiver.Balance)

	store, err := db.GetSQLStore(setup.GetStore())
	require.NoError(t, err)
	require.NotEmpty(t, store)

	n := 5
	// Create a pgtype.Numeric for the amount
	amount := pgtype.Numeric{}
	err = amount.Set(10)
	require.NoError(t, err)

	errs := make(chan error)
	results := make(chan schemas.TransferTxResult)

	for i := 0; i < n; i++ {
		go func(senderID, receiverID int32, transferAmount pgtype.Numeric) {
			result, err := transaction.TransferTx(context.Background(), schemas.TransferTxParams{
				SenderAccountID:   senderID,
				ReceiverAccountID: receiverID,
				Amount:            transferAmount,
			})
			errs <- err
			results <- result
		}(sender.AccountID, receiver.AccountID, amount)
	}

	// Check results
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// Check transaction
		transfer := result.Transaction
		require.NotEmpty(t, transfer)
		require.Equal(t, sender.AccountID, transfer.FromAccountID)
		require.Equal(t, receiver.AccountID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.Equal(t, "COMPLETED", transfer.StatusCode)
		require.NotZero(t, transfer.TransactionID)
		require.NotZero(t, transfer.CreatedAt)

		// Check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, sender.AccountID, fromEntry.AccountID)

		// Create negative amount for comparison
		negativeAmount := pgtype.Numeric{}
		err = negativeAmount.Set(-10)
		require.NoError(t, err)
		require.Equal(t, negativeAmount, fromEntry.Amount)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, receiver.AccountID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)

		// Check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, sender.AccountID, fromAccount.AccountID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, receiver.AccountID, toAccount.AccountID)

		// Get numeric values for calculations
		var senderVal, fromAccountVal, toAccountVal, receiverVal float64
		err = sender.Balance.AssignTo(&senderVal)
		require.NoError(t, err)
		err = fromAccount.Balance.AssignTo(&fromAccountVal)
		require.NoError(t, err)
		err = toAccount.Balance.AssignTo(&toAccountVal)
		require.NoError(t, err)
		err = receiver.Balance.AssignTo(&receiverVal)
		require.NoError(t, err)

		// Calculate differences
		diff1 := pgtype.Numeric{}
		diff2 := pgtype.Numeric{}

		err = diff1.Set(senderVal - fromAccountVal)
		require.NoError(t, err)
		err = diff2.Set(toAccountVal - receiverVal)
		require.NoError(t, err)

		require.Equal(t, diff1, diff2)

		var diff1Val float64
		err = diff1.AssignTo(&diff1Val)
		require.NoError(t, err)
		require.True(t, diff1Val > 0)

		var amountVal float64
		err = amount.AssignTo(&amountVal)
		require.NoError(t, err)

		k := int(diff1Val / amountVal)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// Check final balances
	updatedSender, err := store.GetAccount(context.Background(), sender.AccountID)
	require.NoError(t, err)

	updatedReceiver, err := store.GetAccount(context.Background(), receiver.AccountID)
	require.NoError(t, err)

	// Calculate expected final balances
	var senderVal, receiverVal, amountVal float64
	err = sender.Balance.AssignTo(&senderVal)
	require.NoError(t, err)
	err = receiver.Balance.AssignTo(&receiverVal)
	require.NoError(t, err)
	err = amount.AssignTo(&amountVal)
	require.NoError(t, err)

	expectedSenderBalance := pgtype.Numeric{}
	err = expectedSenderBalance.Set(senderVal - (float64(n) * amountVal))
	require.NoError(t, err)

	expectedReceiverBalance := pgtype.Numeric{}
	err = expectedReceiverBalance.Set(receiverVal + (float64(n) * amountVal))
	require.NoError(t, err)

	require.Equal(t, expectedSenderBalance, updatedSender.Balance)
	require.Equal(t, expectedReceiverBalance, updatedReceiver.Balance)

	fmt.Println(">> after:", updatedSender.Balance, updatedReceiver.Balance)
}
