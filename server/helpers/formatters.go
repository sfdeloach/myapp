package helpers

import "fmt"

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
