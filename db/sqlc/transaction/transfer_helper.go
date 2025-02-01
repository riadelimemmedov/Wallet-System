package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/jackc/pgtype"
	db "github.com/riad/banksystemendtoend/db/sqlc"
	"github.com/riad/banksystemendtoend/util/common"
	"github.com/riad/banksystemendtoend/util/schemas"
)

func createTransferEntries(ctx context.Context, q *db.Queries, arg schemas.TransferTxParams) (fromEntry db.Entry, toEntry db.Entry, err error) {
	debitAmount := NegateNumeric(arg.Amount)

	fromEntry, err = q.CreateEntry(ctx, db.CreateEntryParams{
		AccountID: sql.NullInt32{Int32: arg.SenderAccountID, Valid: true},
		Amount:    debitAmount,
	})

	if err != nil {
		return
	}

	toEntry, err = q.CreateEntry(ctx, db.CreateEntryParams{
		AccountID: sql.NullInt32{Int32: arg.ReceiverAccountID, Valid: true},
		Amount:    arg.Amount,
	})
	if err != nil {
		return
	}
	return fromEntry, toEntry, nil
}

func createTransferTransaction(ctx context.Context, q *db.Queries, arg schemas.TransferTxParams) (db.Transaction, error) {
	transaction, err := q.CreateTransaction(ctx, db.CreateTransactionParams{
		FromAccountID:   sql.NullInt32{Int32: arg.SenderAccountID, Valid: true},
		ToAccountID:     sql.NullInt32{Int32: arg.ReceiverAccountID, Valid: true},
		TypeCode:        "TRANSFER",
		Amount:          arg.Amount,
		CurrencyCode:    arg.CurrencyCode,
		ExchangeRate:    arg.ExchangeRate,
		StatusCode:      "PENDING",
		Description:     sql.NullString{String: "Fund Transfer", Valid: true},
		ReferenceNumber: sql.NullString{String: common.RandomString(10), Valid: true},
		TransactionDate: time.Now(),
	})
	if err != nil {
		return db.Transaction{}, nil
	}
	return transaction, nil
}

func getTransferEntries(ctx context.Context, q *db.Queries, txn db.Transaction, arg schemas.TransferTxParams) (fromEntry db.Entry, toEntry db.Entry, err error) {
	fromEntry, toEntry, err = createTransferEntries(ctx, q, arg)
	return fromEntry, toEntry, err
}

func getTransferTransaction(ctx context.Context, q *db.Queries, arg schemas.TransferTxParams) (db.Transaction, error) {
	transaction, err := createTransferTransaction(ctx, q, arg)
	return transaction, err
}

func updateAccountBalances(ctx context.Context, q *db.Queries, arg schemas.TransferTxParams) (senderAccount db.Account, receiverAccount db.Account, err error) {
	if arg.SenderAccountID < arg.ReceiverAccountID {
		senderAccount, receiverAccount, err = addMoney(ctx, q, arg.SenderAccountID, NegateNumeric(arg.Amount), arg.ReceiverAccountID, arg.Amount)
	} else {
		receiverAccount, senderAccount, err = addMoney(ctx, q, arg.ReceiverAccountID, arg.Amount, arg.SenderAccountID, NegateNumeric(arg.Amount))
	}
	return
}

func addMoney(ctx context.Context, q *db.Queries, senderAccountID int32, senderAmount pgtype.Numeric, receiverAccountID int32, receiverAmount pgtype.Numeric) (senderAccount db.Account, receiverAccount db.Account, err error) {
	senderAccount, err = q.UpdateAccountBalance(ctx, db.UpdateAccountBalanceParams{
		AccountID: senderAccount.AccountID,
		Amount:    senderAmount,
	})
	if err != nil {
		return
	}

	receiverAccount, err = q.UpdateAccountBalance(ctx, db.UpdateAccountBalanceParams{
		AccountID: receiverAccount.AccountID,
		Amount:    senderAmount,
	})
	if err != nil {
		return
	}
	return senderAccount, receiverAccount, nil
}
