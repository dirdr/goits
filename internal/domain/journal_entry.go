package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type JournalEntry struct {
	EntryID       int64           `json:"entry_id"`
	TransactionID uuid.UUID       `json:"transaction_id"`
	AccountID     int64           `json:"account_id"`
	Amount        decimal.Decimal `json:"amount"`
	CreatedAt     time.Time       `json:"created_at"`
}
