package transaction

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgtype"
	db "github.com/riad/banksystemendtoend/db/sqlc"
	"github.com/riad/banksystemendtoend/util/common"
	"github.com/riad/banksystemendtoend/util/config"
	"github.com/riad/banksystemendtoend/util/schemas"
)

func createTransferEntries(ctx context.Context, q *db.Queries, arg schemas.TransferTxParams) (fromEntry db.Entry, toEntry db.Entry, err error) {
	debitAmount := NegateNumeric(arg.Amount)

	fromEntry, err = q.CreateEntry(ctx, db.CreateEntryParams{
		AccountID: sql.NullInt32{Int32: arg.SenderAccountID, Valid: true},
		Amount:    debitAmount,
	})
	if err != nil {
		return fromEntry, toEntry, fmt.Errorf("error creating debit entry: %v", err)
	}

	toEntry, err = q.CreateEntry(ctx, db.CreateEntryParams{
		AccountID: sql.NullInt32{Int32: arg.ReceiverAccountID, Valid: true},
		Amount:    arg.Amount,
	})
	if err != nil {
		return fromEntry, toEntry, fmt.Errorf("error creating credit entry: %v", err)
	}
	return fromEntry, toEntry, nil
}

func createTransferTransaction(ctx context.Context, q *db.Queries, arg schemas.TransferTxParams) (db.Transaction, error) {
	transactionType, _ := CreateTransactionType(config.TransactionTypes.TRANSFER)
	currency, _ := CreateCurrencyCode(config.TransactionCurrencies.USD.CODE)
	transaction_status, _ := CreateTransactionStatus(config.TransactionStatuses.PENDING)

	exchangeRate, err := common.RandomNumeric()
	if err != nil {
		return db.Transaction{}, fmt.Errorf("error creating exchange rate: %v", err)
	}

	transaction, err := q.CreateTransaction(ctx, db.CreateTransactionParams{
		FromAccountID:   sql.NullInt32{Int32: arg.SenderAccountID, Valid: true},
		ToAccountID:     sql.NullInt32{Int32: arg.ReceiverAccountID, Valid: true},
		TypeCode:        transactionType.TypeCode,
		Amount:          arg.Amount,
		CurrencyCode:    currency.CurrencyCode,
		ExchangeRate:    exchangeRate,
		StatusCode:      transaction_status.StatusCode,
		Description:     sql.NullString{String: "Fund Transfer", Valid: true},
		ReferenceNumber: sql.NullString{String: common.RandomString(10), Valid: true},
		TransactionDate: time.Now(),
	})
	if err != nil {
		return db.Transaction{}, fmt.Errorf("error creating transfer transaction: %v", err)
	}
	return transaction, nil
}

func updateAccountBalances(ctx context.Context, q *db.Queries,
	arg schemas.TransferTxParams) (senderAccount db.Account, receiverAccount db.Account, err error) {

	if arg.SenderAccountID < arg.ReceiverAccountID {
		senderAccount, receiverAccount, err = addMoney(ctx, q,
			arg.SenderAccountID, NegateNumeric(arg.Amount),
			arg.ReceiverAccountID, arg.Amount)
	} else {
		receiverAccount, senderAccount, err = addMoney(ctx, q,
			arg.ReceiverAccountID, arg.Amount,
			arg.SenderAccountID, NegateNumeric(arg.Amount))
	}

	if err != nil {
		return senderAccount, receiverAccount, fmt.Errorf("error updating account balances: %v", err)
	}

	return senderAccount, receiverAccount, nil
}

func addMoney(ctx context.Context, q *db.Queries,
	accountID1 int32, amount1 pgtype.Numeric,
	accountID2 int32, amount2 pgtype.Numeric) (account1 db.Account, account2 db.Account, err error) {

	account1, err = q.UpdateAccountBalance(ctx, db.UpdateAccountBalanceParams{
		AccountID: accountID1,
		Amount:    amount1,
	})
	if err != nil {
		return account1, account2, fmt.Errorf("error updating first account balance: %v", err)
	}

	account2, err = q.UpdateAccountBalance(ctx, db.UpdateAccountBalanceParams{
		AccountID: accountID2,
		Amount:    amount2,
	})
	if err != nil {
		return account1, account2, fmt.Errorf("error updating second account balance: %v", err)
	}

	return account1, account2, nil
}
