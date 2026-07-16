// Package models contains the GORM models. Only database concerns belong
// here — business logic lives in services, and models are never sent to the
// frontend directly (DTOs are).
package models

import (
	"time"

	"gorm.io/gorm"
)

// User is the application user record.
type User struct {
	gorm.Model

	Name            string     `gorm:"size:255;not null"`
	Email           string     `gorm:"size:255;not null;uniqueIndex"`
	Password        string     `gorm:"size:255;not null"` // bcrypt hash
	EmailVerifiedAt *time.Time
}
