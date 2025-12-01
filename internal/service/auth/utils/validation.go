package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"regexp"
	"strings"
)

// E.164 phone number format regex (international format)
var phoneRegex = regexp.MustCompile(`^\+[1-9]\d{1,14}$`)

// ValidatePhoneNumber validates phone number in E.164 format
func ValidatePhoneNumber(phone string) error {
	// Trim whitespace
	phone = strings.TrimSpace(phone)

	if phone == "" {
		return fmt.Errorf("phone number cannot be empty")
	}

	if !phoneRegex.MatchString(phone) {
		return fmt.Errorf("invalid phone number format. Expected E.164 format (e.g., +962791234567)")
	}

	return nil
}

// NormalizePhoneNumber normalizes phone number to E.164 format
func NormalizePhoneNumber(phone string) string {
	// Remove spaces, dashes, and parentheses
	phone = strings.TrimSpace(phone)
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.ReplaceAll(phone, "(", "")
	phone = strings.ReplaceAll(phone, ")", "")

	// Add + if missing
	if !strings.HasPrefix(phone, "+") {
		phone = "+" + phone
	}

	return phone
}

// ValidateName validates user name
func ValidateName(name string) error {
	name = strings.TrimSpace(name)

	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	if len(name) < 2 {
		return fmt.Errorf("name must be at least 2 characters long")
	}

	if len(name) > 100 {
		return fmt.Errorf("name must not exceed 100 characters")
	}

	return nil
}

// GenerateOTP generates a random OTP code with specified length
func GenerateOTP(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("OTP length must be positive")
	}

	const digits = "0123456789"
	result := make([]byte, length)

	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", fmt.Errorf("failed to generate OTP: %w", err)
		}
		result[i] = digits[num.Int64()]
	}

	return string(result), nil
}
