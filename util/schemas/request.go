package schemas

import (
	"github.com/jackc/pgtype"
)

type TransferTxParams struct {
	SenderAccountID   int32
	ReceiverAccountID int32
	Amount            pgtype.Numeric
	CurrencyCode      string
	TypeCode          string
	StatusCode        string
	Description       string
	ExchangeRate      pgtype.Numeric
	ReferenceNumber   string
}
