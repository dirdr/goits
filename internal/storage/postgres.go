package storage

import (
	"fmt"
	"log/slog"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	TimeZone string
}

func NewPostgresDB(cfg Config, appLogger *slog.Logger) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode, cfg.TimeZone)

	appLogger.Info("Database DSN", "dsn", dsn)

	gormConfig := &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Warn),
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	appLogger.Info("Running database migrations...")
	err = db.AutoMigrate(&GormAccount{}, &GormTransferEvent{}, &GormJournalEntry{}, &GormAccountBalance{})
	if err != nil {
		return nil, fmt.Errorf("failed to auto migrate database: %w", err)
	}
	appLogger.Info("Database migrations completed.")

	return db, nil
}
