package storage

import (
	"time"

	"github.com/dirdr/goits/internal/domain"
	"github.com/shopspring/decimal"
)

type GormJournalEntry struct {
	EntryID       uint             `gorm:"primaryKey;autoIncrement"`
	TransactionID string           `gorm:"type:varchar(36);not null;index"`
	AccountID     uint             `gorm:"not null;index"`
	Amount        decimal.Decimal  `gorm:"type:numeric(20,8);not null"`
	Type          domain.EntryType `gorm:"type:varchar(50);not null"`
	SourceEventID uint             `gorm:"not null"`
	CreatedAt     time.Time        `gorm:"not null"`
}

func (GormJournalEntry) TableName() string {
	return "journal_entries"
}
