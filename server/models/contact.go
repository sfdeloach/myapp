package models

import "gorm.io/gorm"

type Contact struct {
	gorm.Model        // provides ID, CreatedAt, UpdatedAt, and DeletedAt
	First      string `gorm:"type:varchar(100);not null"`
	Last       string `gorm:"type:varchar(100);not null"`
	Phone      string `gorm:"type:varchar(12)"`
	Email      string `gorm:"type:varchar(255)"`
}
