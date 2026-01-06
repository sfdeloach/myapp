package helpers

import (
	"errors"
	"regexp"
	"server/models"
	"strings"

	"gorm.io/gorm"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// ValidateContact performs custom validation on a Contact and sanitizes input
func ValidateContact(contact *models.Contact, db *gorm.DB) error {
	// Trim leading and trailing white space
	contact.First = strings.TrimSpace(contact.First)
	contact.Last = strings.TrimSpace(contact.Last)
	contact.Phone = strings.TrimSpace(contact.Phone)
	contact.Email = strings.TrimSpace(contact.Email)

	if contact.First == "" {
		return errors.New("first name is required")
	}
	if contact.Last == "" {
		return errors.New("last name is required")
	}

	// Validate and check email for uniqueness
	if err := ValidateEmail(contact.Email, db, contact.ID); err != nil {
		return err
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

// ValidateEmail validates email format and checks for uniqueness in the database for new or updated emails
// The contactID parameter should be 0 for new contacts, or the existing contact's ID for updates
func ValidateEmail(email string, db *gorm.DB, contactID uint) error {
	email = strings.TrimSpace(email)

	// Check format
	if email != "" && !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}

	// Check uniqueness if email is provided and db is available
	if email != "" && db != nil {
		var existingContact models.Contact
		query := db.Where("email = ?", email)

		// Exclude the current contact when updating
		if contactID != 0 {
			query = query.Where("id != ?", contactID)
		}

		if err := query.First(&existingContact).Error; err == nil {
			// Found a contact with this email
			return errors.New("email already in use")
		} else if err != gorm.ErrRecordNotFound {
			// Database error (not a "not found" error)
			return errors.New("database error checking email uniqueness")
		}
	}

	// Returning nil indicates a valid email
	return nil
}
