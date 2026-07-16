// Package seeders populates the database with initial/demo data.
package seeders

import (
	"fmt"
	"log/slog"

	"gorm.io/gorm"

	"flux/app/models"
	"flux/framework"
)

// Run executes every seeder.
func Run(db *gorm.DB) error {
	if err := seedUsers(db); err != nil {
		return err
	}
	slog.Info("seeding complete")
	return nil
}

// seedUsers creates an admin plus a handful of demo users. Idempotent —
// existing emails are skipped.
func seedUsers(db *gorm.DB) error {
	password, err := framework.HashPassword("password")
	if err != nil {
		return err
	}

	users := []models.User{
		{Name: "Admin", Email: "admin@flux.local", Password: password},
	}
	for i := 1; i <= 15; i++ {
		users = append(users, models.User{
			Name:     fmt.Sprintf("Demo User %d", i),
			Email:    fmt.Sprintf("demo%d@flux.local", i),
			Password: password,
		})
	}

	for _, user := range users {
		var count int64
		if err := db.Model(&models.User{}).Where("email = ?", user.Email).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			continue
		}
		if err := db.Create(&user).Error; err != nil {
			return err
		}
	}
	return nil
}
