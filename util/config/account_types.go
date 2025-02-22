package config

// AccountType defines different types of financial accounts
type AccountType struct {
	SAVINGS       string
	CHECKING      string
	FIXED_DEPOSIT string
	MONEY_MARKET  string
}

// Pre-defined account types
var AccountTypes = AccountType{
	SAVINGS:       "SAVINGS",
	CHECKING:      "CHECKING",
	FIXED_DEPOSIT: "FIXED_DEPOSIT",
	MONEY_MARKET:  "MONEY_MARKET",
}
