package database

import (
	"fmt"
	"log"

	"gorm.io/gorm"

	"server/models"
)

// Seed populates the database with test data
func Seed(db *gorm.DB) error {
	log.Println("Starting database seeding...")

	// Clear existing data (optional - remove if you want to keep existing data)
	if err := db.Exec("TRUNCATE TABLE contacts RESTART IDENTITY CASCADE").Error; err != nil {
		return err
	}

	// Generate sample contacts for testing
	numContacts := 8192
	contacts := make([]models.Contact, numContacts)

	for i := range numContacts {
		contacts[i] = models.Contact{
			First: fmt.Sprintf("FirstName%d", i+1),
			Last:  fmt.Sprintf("LastName%d", i+1),
			Phone: fmt.Sprintf("555%07d", i+1),
			Email: fmt.Sprintf("user%d@example.com", i+1),
		}
	}

	// Insert contacts
	if err := db.Create(&contacts).Error; err != nil {
		return err
	}

	log.Printf("Successfully seeded %d contacts", len(contacts))
	return nil
}
