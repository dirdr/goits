package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type EntryType string

const (
	Debit  EntryType = "debit"
	Credit EntryType = "credit"
)

type JournalEntry struct {
	EntryID       uint            `json:"entry_id"`
	TransactionID string          `json:"transaction_id"`
	AccountID     uint            `json:"account_id"`
	Amount        decimal.Decimal `json:"amount"`
	Type          EntryType       `json:"type"`
	SourceEventID uint            `json:"source_event_id"`
	CreatedAt     time.Time       `json:"created_at"`
}

