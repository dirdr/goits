package storage

import (
	"time"

	"github.com/shopspring/decimal"
)

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
