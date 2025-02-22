package config

// IsValidAccountType checks if an account type is valid
func IsValidAccountType(accountType string) bool {
	switch accountType {
	case AccountTypes.SAVINGS,
		AccountTypes.CHECKING,
		AccountTypes.FIXED_DEPOSIT,
		AccountTypes.MONEY_MARKET:
		return true
	}
	return false
}
