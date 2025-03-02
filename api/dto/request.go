package dto

import "mime/multipart"

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

// CreateUserRequest defines the input for creating a new user
type CreateUserRequest struct {
	Username     string                `form:"username" binding:"required"`
	Password     string                `form:"password" binding:"required"`
	Email        string                `form:"email" binding:"required,email"`
	FirstName    string                `form:"first_name"`
	LastName     string                `form:"last_name"`
	PhoneNumber  string                `form:"phone_number"`
	ProfileImage *multipart.FileHeader `form:"profile_image"`
}

// UpdateUserRequest defines the input for updating a user
type UpdateUserRequest struct {
	Username     string                `form:"username"`
	Email        string                `form:"email" binding:"omitempty,email"`
	FirstName    string                `form:"first_name"`
	LastName     string                `form:"last_name"`
	PhoneNumber  string                `form:"phone_number"`
	ProfileImage *multipart.FileHeader `form:"profile_image"`
}

// ChangePasswordRequest defines the input for changing a user's password
type ChangePasswordRequest struct {
	CurrentPassword string `form:"current_password" binding:"required"`
	NewPassword     string `form:"new_password" binding:"required"`
	ConfirmPassword string `form:"confirm_password" binding:"required,eqfield=NewPassword"`
}
