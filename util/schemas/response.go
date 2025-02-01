package schemas

import db "github.com/riad/banksystemendtoend/db/sqlc"

type TransferTxResult struct {
	Transaction db.Transaction
	FromEntry   db.Entry
	ToEntry     db.Entry
	FromAccount db.Account
	ToAccount   db.Account
	Status      db.TransactionStatus
	Type        db.TransactionType
	Currency    db.AccountCurrency
}
