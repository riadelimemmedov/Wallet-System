package db

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/jackc/pgtype"
	db "github.com/riad/banksystemendtoend/db/sqlc"
	"github.com/riad/banksystemendtoend/db/sqlc/transaction"
	"github.com/riad/banksystemendtoend/util/config"
	setup "github.com/riad/banksystemendtoend/util/db"
	"github.com/riad/banksystemendtoend/util/schemas"
	"github.com/stretchr/testify/require"
)

func TestTransfer(t *testing.T) {
	n := 3
	existed := make(map[int]bool)
	var wg sync.WaitGroup
	var err error

	store, err := db.GetSQLStore(setup.GetStore())
	require.NoError(t, err)
	require.NotEmpty(t, store)

	sender := createRandomAccount(t)
	receiver := createRandomAccount(t)
	fmt.Println(">> before:", sender.Balance, receiver.Balance)

	completedStatus, err := transaction.CreateTransactionStatus(config.TransactionStatuses.COMPLETED)
	require.NoError(t, err)

	transferType, err := transaction.CreateTransactionType(config.TransactionTypes.TRANSFER)
	require.NoError(t, err)

	currency, err := transaction.CreateCurrencyCode(config.TransactionCurrencies.USD.CODE)
	require.NoError(t, err)

	amount := pgtype.Numeric{}
	err = amount.Set(10)
	require.NoError(t, err)

	//? Create buffered channels
	errs := make(chan error, n)
	transfer_results := make(chan schemas.TransferTxResult, n)

	//? Execute concurrent transfers
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(senderID, receiverID int32, transferAmount pgtype.Numeric) {
			defer wg.Done()

			transfer, err := transaction.TransferTx(context.Background(), schemas.TransferTxParams{
				SenderAccountID:   senderID,
				ReceiverAccountID: receiverID,
				Amount:            transferAmount,
				CurrencyCode:      currency.CurrencyCode,
				TypeCode:          transferType.TypeCode,
				StatusCode:        completedStatus.StatusCode,
			})

			errs <- err
			transfer_results <- transfer
		}(sender.AccountID, receiver.AccountID, amount)
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-transfer_results
		require.NotEmpty(t, result)

		// Verify transaction details
		transfer := result.Transaction
		require.NotEmpty(t, transfer)
		require.Equal(t, sender.AccountID, transfer.FromAccountID.Int32)
		require.Equal(t, receiver.AccountID, transfer.ToAccountID.Int32)
		// Compare numeric values instead of direct pgtype.Numeric comparison
		var expectedAmount, actualAmount float64
		err = amount.AssignTo(&expectedAmount)
		require.NoError(t, err)
		err = transfer.Amount.AssignTo(&actualAmount)
		require.NoError(t, err)
		require.Equal(t, expectedAmount, actualAmount)
		require.Equal(t, config.TransactionStatuses.COMPLETED, transfer.StatusCode)
		require.NotZero(t, transfer.TransactionID)
		require.NotZero(t, transfer.CreatedAt)

		//? Verify entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, sender.AccountID, fromEntry.AccountID.Int32)

		//? Verify negative amount in from entry
		var fromEntryAmount float64
		err = fromEntry.Amount.AssignTo(&fromEntryAmount)
		require.NoError(t, err)
		require.Equal(t, -10.0, fromEntryAmount)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, receiver.AccountID, toEntry.AccountID.Int32)

		//? Compare the to entry amount
		var toEntryAmount float64
		err = toEntry.Amount.AssignTo(&toEntryAmount)
		require.NoError(t, err)
		require.Equal(t, float64(10), toEntryAmount)

		// Verify accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, sender.AccountID, fromAccount.AccountID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, receiver.AccountID, toAccount.AccountID)

		//? Get numeric values for calculations
		var senderVal, fromAccountVal, toAccountVal, receiverVal float64
		err = sender.Balance.AssignTo(&senderVal)
		require.NoError(t, err)
		err = fromAccount.Balance.AssignTo(&fromAccountVal)
		require.NoError(t, err)
		err = toAccount.Balance.AssignTo(&toAccountVal)
		require.NoError(t, err)
		err = receiver.Balance.AssignTo(&receiverVal)
		require.NoError(t, err)

		//? Calculate and verify differences
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

	wg.Wait()

	//? Verify final account balances
	updatedSender, err := store.GetAccount(context.Background(), sender.AccountID)
	require.NoError(t, err)

	updatedReceiver, err := store.GetAccount(context.Background(), receiver.AccountID)
	require.NoError(t, err)

	//? Calculate expected final balances
	var senderVal, receiverVal, amountVal float64
	err = sender.Balance.AssignTo(&senderVal)
	require.NoError(t, err)
	err = receiver.Balance.AssignTo(&receiverVal)
	require.NoError(t, err)
	err = amount.AssignTo(&amountVal)
	require.NoError(t, err)

	//? Calculate expected final balances
	expectedSenderVal := senderVal - (float64(n) * amountVal)
	expectedReceiverVal := receiverVal + (float64(n) * amountVal)

	//? Get actual final balances
	var updatedSenderVal, updatedReceiverVal float64
	err = updatedSender.Balance.AssignTo(&updatedSenderVal)
	require.NoError(t, err)
	err = updatedReceiver.Balance.AssignTo(&updatedReceiverVal)
	require.NoError(t, err)

	//? Compare the actual values
	require.Equal(t, expectedSenderVal, updatedSenderVal)
	require.Equal(t, expectedReceiverVal, updatedReceiverVal)

	fmt.Println(">> after:", updatedSender.Balance, updatedReceiver.Balance)
	defer CleanupDB(t)
}
