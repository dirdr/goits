package storage

import (
	"time"

	"github.com/shopspring/decimal"
)

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
