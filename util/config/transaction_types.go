package config

// TransactionType defines different types of financial transactions
type TransactionType struct {
	TRANSFER   string
	DEPOSIT    string
	WITHDRAWAL string
	PAYMENT    string
	REFUND     string
	ADJUSTMENT string
	FEE        string
	INTEREST   string
}

// TransactionStatus defines possible states of a transaction
type TransactionStatus struct {
	PENDING   string
	COMPLETED string
	FAILED    string
	CANCELLED string
	REVERSED  string
}

// Currency defines possible states of a currencies
type Currency struct {
	CODE   string
	NAME   string
	SYMBOL string
	RATE   float64
}

// Pre-defined transaction types
var TransactionTypes = TransactionType{
	TRANSFER:   "TRANSFER",
	DEPOSIT:    "DEPOSIT",
	WITHDRAWAL: "WITHDRAWAL",
	PAYMENT:    "PAYMENT",
	REFUND:     "REFUND",
	ADJUSTMENT: "ADJUSTMENT",
	FEE:        "FEE",
	INTEREST:   "INTEREST",
}

// Pre-defined transaction statuses
var TransactionStatuses = TransactionStatus{
	PENDING:   "PENDING",
	COMPLETED: "COMPLETED",
	FAILED:    "FAILED",
	CANCELLED: "CANCELLED",
	REVERSED:  "REVERSED",
}

// Pre-defined currencies statuses
var TransactionCurrencies = struct {
	USD Currency
	EUR Currency
	GBP Currency
	JPY Currency
}{
	USD: Currency{
		CODE:   "USD",
		NAME:   "US Dollar",
		SYMBOL: "$",
		RATE:   1.0,
	},
	EUR: Currency{
		CODE:   "EUR",
		NAME:   "Euro",
		SYMBOL: "€",
		RATE:   1.08,
	},
	GBP: Currency{
		CODE:   "GBP",
		NAME:   "British Pound",
		SYMBOL: "£",
		RATE:   1.26,
	},
	JPY: Currency{
		CODE:   "Japanese Yen",
		NAME:   "Japanese Yen",
		SYMBOL: "¥",
		RATE:   0.0067,
	},
}
