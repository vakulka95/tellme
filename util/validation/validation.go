package validation

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

//
// Validation: s, maximum: no limit, at least 1 digit, at least 1 upper case
//
func ValidatePassword(password string) error {
	//if !hasAtLeastNNumbers(password, 1) {
	//	return fmt.Errorf("password has no digits")
	//}
	//
	//if !hasAtLeastNUpperCases(password, 1) {
	//	return fmt.Errorf("password has no uppercases")
	//}
	//
	//if !hasAtLeastNLowerCases(password, 1) {
	//	return fmt.Errorf("password has no lowercases")
	//}

	if !hasAtLeastNLength(password, 7) {
		return fmt.Errorf("password has length < 10")
	}

	return nil
}

//
// Validation: if @ symbol is in the text string, consider input valid
//
func ValidateEmail(email string) error {
	if !hasAtLeastNSymbols(email, 1, "@") {
		return fmt.Errorf("email has invalid format")
	}

	return nil
}

//
// Must be no less than (≥) 10 digits
//
func ValidatePhoneNumber(phoneNumber string) error {
	if !regexp.MustCompile(`^\d{10,}$`).MatchString(phoneNumber) {
		return fmt.Errorf("phone number is invalid")
	}

	return nil
}

func hasAtLeastNNumbers(s string, n int) bool {
	for i := range s {
		if unicode.IsNumber(rune(s[i])) {
			return true
		}
	}

	return false
}

func hasAtLeastNUpperCases(s string, n int) bool {
	for i := range s {
		if unicode.IsUpper(rune(s[i])) {
			return true
		}
	}

	return false
}

func hasAtLeastNLowerCases(s string, n int) bool {
	for i := range s {
		if unicode.IsLower(rune(s[i])) {
			return true
		}
	}

	return false
}

func hasAtLeastNLength(s string, n int) bool {
	return len(s) >= n
}

func hasAtLeastNSymbols(s string, n int, symbol string) bool {
	return strings.Count(s, symbol) >= n
}
