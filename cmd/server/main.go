package main

import (
	"log/slog"

	"github.com/dirdr/goits/internal/config"
	"github.com/dirdr/goits/internal/handler"
	"github.com/dirdr/goits/internal/service"
	"github.com/dirdr/goits/internal/storage"
	"github.com/dirdr/goits/pkg/logger"

	_ "github.com/dirdr/goits/docs"
)

func main() {
	appLogger := initLogger()

	cfg, err := config.LoadConfig()
	if err != nil {
		appLogger.Error("Failed to load configuration", "error", err)
		return
	}

	dbConfig := storage.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
		TimeZone: cfg.Database.TimeZone,
	}

	db, err := storage.NewPostgresDB(dbConfig, appLogger)
	if err != nil {
		appLogger.Error("Failed to connect to database", "error", err)
		return
	}

	accountRepo := storage.NewGormAccountRepository(db)
	accountBalanceRepo := storage.NewGormAccountBalanceRepository(db)
	transferEventRepo := storage.NewGormTransferEventRepository(db)
	journalRepo := storage.NewGormJournalRepository(db)

	accountService := service.NewAccountService(accountRepo, accountBalanceRepo)
	transactionService := service.NewTransactionService(accountRepo, accountBalanceRepo, transferEventRepo, journalRepo)
	integrityService := service.NewIntegrityService(journalRepo)

	r := handler.GetRouter(accountService, transactionService, integrityService, appLogger, db)

	appLogger.Info("Server starting", "port", cfg.Server.Port)
	if err := r.Run(cfg.Server.Port); err != nil {
		appLogger.Error("Failed to start server", "error", err)
	}
}

func initLogger() *slog.Logger {
	logger := logger.New("info")
	slog.SetDefault(logger)

	logger.Info("Starting goits application...")

	return logger
}
