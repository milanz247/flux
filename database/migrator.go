package database

import (
	"fmt"
	"log/slog"
	"sort"
	"time"

	"gorm.io/gorm"
)

// Migration is a single schema change. Migrations register themselves in an
// init() function (see database/migrations) and run in ID order; applied IDs
// are tracked in the `migrations` table so each runs exactly once.
type Migration struct {
	ID string // e.g. "0001_create_users_table"
	Up func(*gorm.DB) error
}

var registry []Migration

// RegisterMigration adds a migration to the global registry. Called from
// init() in each migration file.
func RegisterMigration(m Migration) {
	registry = append(registry, m)
}

type migrationRecord struct {
	ID        uint   `gorm:"primarykey"`
	Migration string `gorm:"size:255;uniqueIndex"`
	AppliedAt time.Time
}

func (migrationRecord) TableName() string { return "migrations" }

// Migrate runs all pending migrations.
func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&migrationRecord{}); err != nil {
		return fmt.Errorf("database: creating migrations table: %w", err)
	}

	sort.Slice(registry, func(i, j int) bool { return registry[i].ID < registry[j].ID })

	applied := map[string]bool{}
	var records []migrationRecord
	if err := db.Find(&records).Error; err != nil {
		return err
	}
	for _, record := range records {
		applied[record.Migration] = true
	}

	ran := 0
	for _, migration := range registry {
		if applied[migration.ID] {
			continue
		}
		slog.Info("migrating", slog.String("migration", migration.ID))
		if err := migration.Up(db); err != nil {
			return fmt.Errorf("database: migration %s: %w", migration.ID, err)
		}
		record := migrationRecord{Migration: migration.ID, AppliedAt: time.Now()}
		if err := db.Create(&record).Error; err != nil {
			return err
		}
		ran++
	}

	if ran == 0 {
		slog.Info("nothing to migrate")
	} else {
		slog.Info("migrations complete", slog.Int("ran", ran))
	}
	return nil
}
