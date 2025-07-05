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
	log := logger.New("info")
	slog.SetDefault(log)

	log.Info("Starting goits application...")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error("Failed to load configuration", "error", err)
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

	log.Info(dbConfig.Host)
	log.Info(dbConfig.Port)
	log.Info(dbConfig.User)
	log.Info(dbConfig.Password)
	log.Info(dbConfig.DBName)
	log.Info(dbConfig.SSLMode)
	log.Info(dbConfig.TimeZone)

	db, err := storage.NewPostgresDB(dbConfig, log)
	if err != nil {
		log.Error("Failed to connect to database", "error", err)
		return
	}

	accountRepo := storage.NewGormAccountRepository(db)
	accountBalanceRepo := storage.NewGormAccountBalanceRepository(db)
	transferEventRepo := storage.NewGormTransferEventRepository(db)
	journalRepo := storage.NewGormJournalRepository(db)

	accountService := service.NewAccountService(accountRepo, accountBalanceRepo, db)
	transactionService := service.NewTransactionService(accountRepo, accountBalanceRepo, transferEventRepo, journalRepo, db)

	r := handler.GetRouter(accountService, transactionService, log)

	log.Info("Server starting", "port", cfg.Server.Port)
	if err := r.Run(cfg.Server.Port); err != nil {
		log.Error("Failed to start server", "error", err)
	}
}
