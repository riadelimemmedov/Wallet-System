package schemas

import (
	"github.com/jackc/pgtype"
)

type TransferTxParams struct {
	SenderAccountID   int32
	ReceiverAccountID int32
	Amount            pgtype.Numeric
	CurrencyCode      string
	ExchangeRate      pgtype.Numeric
}
