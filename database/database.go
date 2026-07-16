// Package database opens the MySQL connection and runs migrations/seeders.
package database

import (
	"fmt"
	"log/slog"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"flux/config"
)

// EnsureDatabase creates the configured database if it does not exist yet
// (used by `flux migrate` so a fresh checkout works without manual SQL).
func EnsureDatabase(cfg *config.Config) error {
	serverDSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=%s",
		cfg.Database.User, cfg.Database.Password,
		cfg.Database.Host, cfg.Database.Port, cfg.Database.Charset)

	db, err := gorm.Open(mysql.Open(serverDSN), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		return fmt.Errorf("database: connecting to mysql server: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	create := fmt.Sprintf(
		"CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET %s COLLATE %s_unicode_ci",
		cfg.Database.Name, cfg.Database.Charset, cfg.Database.Charset)
	return db.Exec(create).Error
}

// Connect opens the GORM MySQL connection with sane pool defaults.
func Connect(cfg *config.Config) (*gorm.DB, error) {
	logLevel := gormlogger.Warn
	if cfg.App.Debug {
		logLevel = gormlogger.Info
	}

	db, err := gorm.Open(mysql.Open(cfg.Database.DSN()), &gorm.Config{
		Logger: gormlogger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("database: connecting to mysql: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	slog.Info("database connected",
		slog.String("host", cfg.Database.Host),
		slog.String("database", cfg.Database.Name),
	)
	return db, nil
}
