package helpers

import (
	"errors"
	"regexp"
	"server/models"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// ValidateContact performs custom validation on a Contact and sanitizes input
func ValidateContact(contact *models.Contact) error {
	// Sanitize and validate required fields
	contact.First = strings.TrimSpace(contact.First)
	contact.Last = strings.TrimSpace(contact.Last)

	if contact.First == "" {
		return errors.New("first name is required")
	}
	if contact.Last == "" {
		return errors.New("last name is required")
	}

	// Sanitize and validate email
	contact.Email = strings.TrimSpace(contact.Email)
	if contact.Email != "" && !emailRegex.MatchString(contact.Email) {
		return errors.New("invalid email format")
	}

	// Sanitize phone: extract only digits
	if contact.Phone != "" {
		digits := ""
		for _, ch := range contact.Phone {
			if ch >= '0' && ch <= '9' {
				digits += string(ch)
			}
		}

		// Validate phone format (10 or 11 digits)
		if len(digits) != 10 && len(digits) != 11 {
			return errors.New("phone must be 10 or 11 digits")
		}

		// Store only the digits
		contact.Phone = digits
	}

	return nil
}
