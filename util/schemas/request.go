package schemas

import (
	"database/sql"
	"time"

	"github.com/jackc/pgtype"
)

type CreateTransactionParams struct {
	FromAccountID   sql.NullInt32  `json:"from_account_id"`
	ToAccountID     sql.NullInt32  `json:"to_account_id"`
	TypeCode        string         `json:"type_code"`
	Amount          pgtype.Numeric `json:"amount"`
	CurrencyCode    string         `json:"currency_code"`
	ExchangeRate    pgtype.Numeric `json:"exchange_rate"`
	StatusCode      string         `json:"status_code"`
	Description     sql.NullString `json:"description"`
	ReferenceNumber sql.NullString `json:"reference_number"`
	TransactionDate time.Time      `json:"transaction_date"`
}
