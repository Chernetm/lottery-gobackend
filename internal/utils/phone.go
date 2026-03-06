package utils

import (
	"errors"
	"strings"
)

// ValidateAndFormatPhoneNumber validates that the phone number starts with 0 and is 10 digits.
// It then formats it to international format (starting with 251 instead of 0).
func ValidateAndFormatPhoneNumber(phone string) (string, error) {
	// 1. Remove any whitespace
	phone = strings.TrimSpace(phone)

	// 2. Check length (must be 10)
	if len(phone) != 10 {
		return "", errors.New("phone number must be exactly 10 digits")
	}

	// 3. Check if starts with '0'
	if !strings.HasPrefix(phone, "0") {
		return "", errors.New("phone number must start with 0")
	}

	// 4. Check if all are digits
	for _, r := range phone {
		if r < '0' || r > '9' {
			return "", errors.New("phone number must contain only digits")
		}
	}

	// 5. Format: replace '0' with '251'
	formatted := "251" + phone[1:]

	return formatted, nil
}
