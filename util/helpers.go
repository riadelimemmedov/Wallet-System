package util

import (
	"fmt"

	"github.com/jackc/pgtype"
)

func SetNumeric(value string) (pgtype.Numeric, error) {
	var numeric pgtype.Numeric
	err := numeric.Set(value)
	if err != nil {
		return numeric, fmt.Errorf("failed to set numeric value for exchange rate")
	}
	return numeric, nil
}
