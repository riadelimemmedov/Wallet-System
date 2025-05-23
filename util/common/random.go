package common

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unicode"

	"github.com/jackc/pgtype"
	"github.com/shopspring/decimal"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

var (
	minMoney int64 = 0
	maxMoney int64 = 1000000
)

// const numbers = "0123456789"

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

// ! RandomInt generates a random integer between min and max number
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// ! RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	random_string := sb.String()
	if len(random_string) > 0 {
		first_char := unicode.ToUpper(rune(random_string[0]))
		return string(first_char) + random_string[1:]
	}
	return random_string
}

// !RandomUsername generates a random username
func RandomUsername() string {
	return fmt.Sprintf("%s_%d", RandomString(6), RandomInt(100, 999))
}

// !RandomEmail generates a random email
func RandomEmail() string {
	return fmt.Sprintf("%s@gmail.com", RandomString(10))
}

// !RandomFirstName generates a random first name
func RandomFirstName() string {
	names := []string{
		"John", "Jane", "Michael", "Emma", "William",
		"Olivia", "James", "Sophia", "Robert", "Isabella",
	}
	return names[rand.Intn(len(names))]
}

// ! RandomLastName generates a random last name
func RandomLastName() string {
	names := []string{
		"Smith", "Johnson", "Williams", "Brown", "Jones",
		"Garcia", "Miller", "Davis", "Rodriguez", "Martinez",
	}
	return names[rand.Intn(len(names))]
}

// ! RandomPhoneNumber generates a random phone number
func RandomPhoneNumber() string {
	return fmt.Sprintf("+1%d%d%d",
		RandomInt(100, 999),
		RandomInt(100, 999),
		RandomInt(1000, 9999),
	)
}

// ! RandomPassword generates a random password
func RandomPassword() string {
	return RandomString(12)
}

// ! RandomBool generates a random boolean
func RandomBool() bool {
	return rand.Intn(2) == 1
}

// ! RandomTime generates a random time within the last n days
func RandomTime(days int) time.Time {
	duration := time.Duration(rand.Intn(days)) * 24 * time.Hour
	return time.Now().Add(-duration)
}

// ! RandomProfileImage generates a random profile image URL
func RandomProfileImage() string {
	imageServices := []string{
		"https://picsum.photos/%d/%d",
		"https://via.placeholder.com/%dx%d",
		"https://placehold.co/%dx%d",
	}

	width := RandomInt(200, 400)
	height := RandomInt(200, 400)

	serviceURL := imageServices[rand.Intn(len(imageServices))]
	return fmt.Sprintf(serviceURL, width, height)

}

// ! RandomAccountNumber generates a random account number
func RandomAccountNumber() string {
	return fmt.Sprintf("ACC%d", RandomInt(100000, 999999))
}

// ! RandomAccountType generates a random account type
func RandomAccountType() string {
	accountTypes := []string{"SAVINGS", "CHECKING", "INVESTMENT", "CREDIT", "FIXED_DEPOSIT", "MONEY_MARKET"}
	// accountTypes := []string{"MONEY_MARKET"}
	return accountTypes[rand.Intn(len(accountTypes))]
}

// ! RandomFloat generates a random float64 between min and max
func RandomFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// ! RandomCurrency generates a random currency code
func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "GBP", "JPY", "CAD", "TL", "AZN", "LR"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

// ! RandomCurrencyName returns the full name of the given currency code
func RandomCurrencyName(code string) string {
	currencyMap := map[string]string{
		"USD": "United States Dollar",
		"EUR": "Euro",
		"GBP": "British Pound Sterling",
		"JPY": "Japanese Yen",
		"CAD": "Canadian Dollar",
		"TL":  "Turkish Lira",
		"AZN": "Azerbaijani Manat",
		"LR":  "Liberian Dollar",
	}

	if name, ok := currencyMap[code]; ok {
		return name
	}
	return "Unknown Currency"
}

// ! RandomCurrencySymbol returns the symbol of the given currency code
func RandomCurrencySymbol(code string) string {
	symbolMap := map[string]string{
		"USD": "$",
		"EUR": "€",
		"GBP": "£",
		"JPY": "¥",
		"CAD": "C$",
		"TL":  "₺",
		"AZN": "₼",
		"LR":  "$",
	}
	if symbol, ok := symbolMap[code]; ok {
		return symbol
	}
	return "?"
}

// ! RandomMoney generates a random amount of money between minMoney and maxMoney
func RandomMoney() decimal.Decimal {
	amount := decimal.New(rand.Int63n(maxMoney), -2)
	return amount
}

// !RandomNumeric generates a random numeric value between 0-10000 with 2 decimal places
func RandomNumeric() (pgtype.Numeric, error) {
	return SetNumeric(fmt.Sprintf("%.2f", RandomFloat(0, 10000)))
}

// ! RandomNegativeNumeric generates a random numeric value between -10000 and 0 with 2 decimal places
func RandomNegativeNumeric() (pgtype.Numeric, error) {
	return SetNumeric(fmt.Sprintf("%.2f", RandomFloat(-10000, 0)))
}

// ! RandomInterestRate generates a random interest rate between 1.1%-99.99%
func RandomInterestRate() (pgtype.Numeric, error) {
	rate := RandomFloat(1.1, 99.99)
	return SetNumeric(fmt.Sprintf("%.2f", rate))
}

// ! RandomTransactionType generates a random transaction type from predefined types
func RandomTransactionType() string {
	transactionTypes := []string{
		"TRANSFER",
		"DEPOSIT",
		"WITHDRAWAL",
		"PAYMENT",
		"REFUND",
		"FEE",
		"INTEREST",
	}
	return transactionTypes[rand.Intn(len(transactionTypes))]
}

// ! RandomTransactionStatus generates a random transaction status from predefined statuses
func RandomTransactionStatus() string {
	transactionStatuses := []string{
		"PENDING",
		"COMPLETED",
		"FAILED",
		"CANCELLED",
		"REJECTED",
		"REVERSED",
		"PROCESSING",
	}
	return transactionStatuses[rand.Intn(len(transactionStatuses))]
}
