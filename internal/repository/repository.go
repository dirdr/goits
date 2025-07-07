package repository

import (
	"context"

	"github.com/dirdr/goits/internal/domain"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type AccountRepository interface {
	CreateAccount(ctx context.Context, tx *gorm.DB, account *domain.Account) error
	GetAccountByID(ctx context.Context, tx *gorm.DB, accountID uint) (*domain.Account, error)
	AccountExists(ctx context.Context, tx *gorm.DB, accountID uint) (bool, error)
}

type AccountBalanceRepository interface {
	GetAccountBalance(ctx context.Context, tx *gorm.DB, accountID uint) (*domain.AccountBalance, error)
	UpsertAccountBalance(ctx context.Context, tx *gorm.DB, balance *domain.AccountBalance) error
	UpdateAccountBalanceWithVersion(ctx context.Context, tx *gorm.DB, balance *domain.AccountBalance, expectedVersion int) error
}

type TransferEventRepository interface {
	SaveTransferEvent(ctx context.Context, tx *gorm.DB, event *domain.TransferEvent) error
}

type JournalRepository interface {
	SaveJournalEntry(ctx context.Context, tx *gorm.DB, entry *domain.JournalEntry) error
	GetTotalsByEntryType(ctx context.Context, tx *gorm.DB) (map[domain.EntryType]decimal.Decimal, error)
}
