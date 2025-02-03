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
