package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal" // The standard for handling money in Go
)

// Account represents a financial account in its purest form.
type Account struct {
	AccountID int64     `json:"account_id"`
	CreatedAt time.Time `json:"created_at"`
}

// Event represents an immutable action that has occurred in the system.
type Event struct {
	EventID   int64  `json:"event_id"`
	EventType string `json:"event_type"`
	// Using json.RawMessage is efficient. It avoids unmarshaling the payload
	// until it's actually needed by a specific event handler.
	Payload   json.RawMessage `json:"payload"`
	CreatedAt time.Time       `json:"created_at"`
}

// JournalEntry represents a single debit or credit in the double-entry system.
// This is the fundamental building block of the ledger.
type JournalEntry struct {
	EntryID       int64     `json:"entry_id"`
	TransactionID uuid.UUID `json:"transaction_id"`
	AccountID     int64     `json:"account_id"`
	// Using decimal.Decimal is CRITICAL for financial calculations to avoid
	// floating-point precision errors that float64 would introduce.
	Amount    decimal.Decimal `json:"amount"`
	CreatedAt time.Time       `json:"created_at"`
}
