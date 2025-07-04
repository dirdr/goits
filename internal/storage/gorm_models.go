package storage

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	AccountID int64 `gorm:"primaryKey;autoIncrement:false"`
	CreatedAt time.Time
}

type Event struct {
	gorm.Model
	EventID   int64 `gorm:"primaryKey;autoIncrement:true"`
	EventType string
	Payload   []byte
	CreatedAt time.Time
}

type JournalEntry struct {
	gorm.Model
	EntryID       int64     `gorm:"primaryKey;autoIncrement:true"`
	TransactionID uuid.UUID `gorm:"type:uuid;index"`
	AccountID     int64     `gorm:"index"`
	Account       Account
	Amount        decimal.Decimal `gorm:"type:numeric(19,4)"`
	CreatedAt     time.Time
}
