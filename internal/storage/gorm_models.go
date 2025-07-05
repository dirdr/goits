package storage

import (
	"time"

	"github.com/dirdr/goits/internal/domain"
	"github.com/shopspring/decimal"
)

// 0. Account (basic entity)
type GormAccount struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (GormAccount) TableName() string {
	return "accounts"
}

// 1. Transfer Events (source of truth)
type GormTransferEvent struct {
	EventID       uint            `gorm:"primaryKey;autoIncrement"`
	TransferID    string          `gorm:"type:varchar(36);not null"`
	FromAccountID uint            `gorm:"not null"`
	ToAccountID   uint            `gorm:"not null"`
	Amount        decimal.Decimal `gorm:"type:numeric(20,8);not null"`
	EventType     string          `gorm:"type:varchar(100);not null"`
	CreatedAt     time.Time       `gorm:"not null"`
}

func (GormTransferEvent) TableName() string {
	return "transfer_events"
}

// 2. Journal Entries (double-entry projection)
type GormJournalEntry struct {
	EntryID       uint             `gorm:"primaryKey;autoIncrement"`
	TransactionID string           `gorm:"type:varchar(36);not null;index"`
	AccountID     uint             `gorm:"not null;index"`
	Amount        decimal.Decimal  `gorm:"type:numeric(20,8);not null"`
	Type          domain.EntryType `gorm:"type:varchar(50);not null"` // "debit" or "credit"
	SourceEventID uint             `gorm:"not null"`                  // Links to event
	CreatedAt     time.Time        `gorm:"not null"`
}

func (GormJournalEntry) TableName() string {
	return "journal_entries"
}

// 3. Account Balances (performance projection)
type GormAccountBalance struct {
	AccountID   uint            `gorm:"primaryKey"`
	Balance     decimal.Decimal `gorm:"type:numeric(20,8);not null"`
	Version     int             `gorm:"not null"`
	LastEventID uint            `gorm:"not null"`
	UpdatedAt   time.Time       `gorm:"not null"`
}

func (GormAccountBalance) TableName() string {
	return "account_balances"
}
