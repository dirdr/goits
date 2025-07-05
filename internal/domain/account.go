package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type Account struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AccountBalance struct {
	AccountID   uint            `json:"account_id"`
	Balance     decimal.Decimal `json:"balance"`
	Version     int             `json:"version"`
	LastEventID uint            `json:"last_event_id"`
	UpdatedAt   time.Time       `json:"updated_at"`
}
