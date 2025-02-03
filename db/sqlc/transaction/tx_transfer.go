package db

import (
	"context"
	"fmt"

	db "github.com/riad/banksystemendtoend/db/sqlc"
	setup "github.com/riad/banksystemendtoend/util/db"
	"github.com/riad/banksystemendtoend/util/schemas"
)

// TransferTx executes a money transfer transaction between two accounts// TransferTx executes a money transfer transaction between two accounts
func TransferTx(ctx context.Context, arg schemas.TransferTxParams) (schemas.TransferTxResult, error) {
	var result schemas.TransferTxResult

	store, err := db.GetSQLStore(setup.GetStore())

	if err != nil {
		return result, fmt.Errorf("failed to get SQL store: %v", err)
	} else {
		err = store.ExecTx(ctx, func(q *db.Queries) error {
			var err error
			//? Step 1: Create Transaction
			result.Transaction, err = createTransferTransaction(ctx, q, arg)
			if err != nil {
				return fmt.Errorf("failed to create transfer transaction: %v", err)
			}

			//? Step 2: Create Entries
			result.FromEntry, result.ToEntry, err = createTransferEntries(ctx, q, arg)
			if err != nil {
				return fmt.Errorf("failed to create transfer entries: %v", err)
			}

			//? Step 3: Update Account Balances
			result.FromAccount, result.ToAccount, err = updateAccountBalances(ctx, q, arg)
			if err != nil {
				return fmt.Errorf("failed to update account balances: %v", err)
			}

			//? Step 4: Update Transaction Status
			result.Transaction, err = q.UpdateTransactionStatus(ctx, db.UpdateTransactionStatusParams{
				TransactionID: result.Transaction.TransactionID,
				StatusCode:    "COMPLETED",
			})
			if err != nil {
				return fmt.Errorf("failed to update transaction status: %v", err)
			}
			return nil
		})

		if err != nil {
			return result, fmt.Errorf("transfer transaction failed: %v", err)
		}

		return result, nil
	}

}
