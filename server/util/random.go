// server/util/random.go

package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

var seededRand *rand.Rand

func init() {
	// Seed the local random number generator with the current time
	seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// RandomInt generates a random integer between a specified minimum and maximum value (inclusive)
func RandomInt(min, max int64) int64 {
	return min + seededRand.Int63n(max-min+1)
}

// RandomString generates a random string of a specified length
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)
	for i := 0; i < n; i++ {
		c := alphabet[seededRand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

// RandomOwner generates a random owner name (assuming a simple 6-character string)
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney generates a random amount of money between 0 and 1000 (inclusive)
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// RandomCurrency generates a random currency from a pre-defined list
func RandomCurrency() string {
	currencies := []string{"EUR", "USD", "CAD"}
	n := len(currencies)
	return currencies[seededRand.Intn(n)]
}

// RandomEmail generates a random email address with a 6-character username and gmail.com domain
func RandomEmail() string {
	return fmt.Sprintf("%s@gmail.com", RandomString(6))
}

// RandomUserStruct creates a random user struct for testing.
func RandomUserStruct() struct {
	Username       string
	FullName       string
	Email          string
	HashedPassword string
	Role           string
} {
	// Generate a random plain password
	plainPassword := RandomString(10)

	// Hash it using bcrypt
	hashed, err := HashPassword(plainPassword)
	if err != nil {
		// In tests, you can panic here safely, since it’s a setup helper
		panic(fmt.Sprintf("failed to hash password: %v", err))
	}

	return struct {
		Username       string
		FullName       string
		Email          string
		HashedPassword string
		Role           string
	}{
		Username:       RandomString(8),
		FullName:       fmt.Sprintf("%s %s", cases.Title(language.English).String(RandomString(5)), cases.Title(language.English).String(RandomString(6))),
		Email:          RandomEmail(),
		HashedPassword: hashed,
		Role:           "student", // You can switch to depositor/banker as needed
	}
}
