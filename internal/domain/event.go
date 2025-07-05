package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type TransferEvent struct {
	EventID       uint            `json:"event_id"`
	TransferID    string          `json:"transfer_id"`
	FromAccountID uint            `json:"from_account_id"`
	ToAccountID   uint            `json:"to_account_id"`
	Amount        decimal.Decimal `json:"amount"`
	EventType     string          `json:"event_type"`
	CreatedAt     time.Time       `json:"created_at"`
}
