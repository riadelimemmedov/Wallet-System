package db

import (
	"math/big"

	"github.com/jackc/pgtype"
)

// NegateNumeric creates a new pgtype.Numeric with negated value while preserving other properties
func NegateNumeric(num pgtype.Numeric) pgtype.Numeric {
	negated := num
	negated.Int = new(big.Int).Neg(num.Int)
	return negated
}
