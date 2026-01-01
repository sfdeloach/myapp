package database

import (
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

	// Sample contacts for testing
	contacts := []models.Contact{
		{
			First: "John",
			Last:  "Doe",
			Phone: "5551234567",
			Email: "john.doe@example.com",
		},
		{
			First: "Jane",
			Last:  "Smith",
			Phone: "5559876543",
			Email: "jane.smith@example.com",
		},
		{
			First: "Bob",
			Last:  "Johnson",
			Phone: "15551112222",
			Email: "bob.johnson@example.com",
		},
		{
			First: "Alice",
			Last:  "Williams",
			Phone: "5554445555",
			Email: "alice.williams@example.com",
		},
		{
			First: "Charlie",
			Last:  "Brown",
			Phone: "5556667777",
			Email: "charlie.brown@example.com",
		},
		{
			First: "Diana",
			Last:  "Davis",
			Phone: "5558889999",
			Email: "diana.davis@example.com",
		},
		{
			First: "Eve",
			Last:  "Martinez",
			Phone: "5552223333",
			Email: "eve.martinez@example.com",
		},
		{
			First: "Frank",
			Last:  "Garcia",
			Phone: "15557778888",
			Email: "frank.garcia@example.com",
		},
	}

	// Insert contacts
	if err := db.Create(&contacts).Error; err != nil {
		return err
	}

	log.Printf("Successfully seeded %d contacts", len(contacts))
	return nil
}
