package helpers

import (
	"fmt"
	"slices"
	"strings"
)

// FormatPhone formats a phone number for display
func FormatPhone(phone string) string {
	// Remove any non-digit characters
	digits := ""
	for _, ch := range phone {
		if ch >= '0' && ch <= '9' {
			digits += string(ch)
		}
	}

	// Format as (XXX) XXX-XXXX for 10-digit numbers
	if len(digits) == 10 {
		return fmt.Sprintf("(%s) %s-%s", digits[0:3], digits[3:6], digits[6:10])
	}

	// Format as X-XXX-XXX-XXXX for 11-digit numbers
	if len(digits) == 11 {
		return fmt.Sprintf("%s-%s-%s-%s", digits[0:1], digits[1:4], digits[4:7], digits[7:11])
	}

	// Return original if not 10 digits
	return phone
}

// GetPage calculates the next or previous page with wrap-around logic
// Assumes page numbers are 1-indexed (pages 1 through totalPages).
func LoadMore(currentPage int) int {
	return currentPage + 1
}

// Receives an integer and returns a string with commas every three digits
func Commas(number string) string {
	// Handle negative numbers by preserving the sign
	isNegative := strings.HasPrefix(number, "-")
	if isNegative {
		number = number[1:]
	}

	var result string
	split := strings.Split(number, "")
	slices.Reverse(split)

	for index, value := range split {
		// Add comma before this digit if it's a multiple of 3 and not the last digit
		if index > 0 && index%3 == 0 {
			result = "," + result
		}
		result = value + result
	}

	// Prepend negative sign if needed
	if isNegative {
		result = "-" + result
	}

	return result
}
