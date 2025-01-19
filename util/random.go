package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unicode"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

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
