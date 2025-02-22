package dto

// CreateUserAccountResponse represents the combined response after creating both user and account
type CreateUserAccountRequest struct {
	// User details
	Username        string `json:"username" binding:"required"`
	Password        string `json:"password" binding:"required,min=6"`
	Email           string `json:"email" binding:"required,email"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	PhoneNumber     string `json:"phone_number"`
	ProfileImageUrl string `json:"profile_image_url"`

	// Account details
	AccountType    string  `json:"account_type" binding:"required"`
	CurrencyCode   string  `json:"currency_code" binding:"required"`
	InterestRate   float64 `json:"interest_rate" binding:"required"`
	OverdraftLimit float64 `json:"overdraft_limit" binding:"required"`
}

// CreateAccountTypeRequest represents the request for creating an account type
type CreateAccountTypeRequest struct {
	AccountType string `json:"account_type" binding:"required,min=2,max=50"`
	Description string `json:"description" binding:"required,min=2,max=200"`
}
