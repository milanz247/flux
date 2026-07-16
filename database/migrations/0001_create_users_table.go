// Package migrations contains the schema migrations. Each file registers
// itself with the migrator in init(); `flux migrate` runs the pending ones.
package migrations

import (
	"gorm.io/gorm"

	"flux/app/models"
	"flux/database"
)

func init() {
	database.RegisterMigration(database.Migration{
		ID: "0001_create_users_table",
		Up: func(db *gorm.DB) error {
			return db.AutoMigrate(&models.User{})
		},
	})
}
