package handler

import (
	"time"

	"github.com/shopspring/decimal"
)

type CreateAccountRequest struct {
	AccountID      uint            `json:"account_id"`
	InitialBalance decimal.Decimal `json:"initial_balance"`
}

type GetAccountResponse struct {
	AccountID uint            `json:"account_id"`
	Balance   decimal.Decimal `json:"balance"`
	Version   int             `json:"version"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type CreateTransactionRequest struct {
	SourceAccountID      uint            `json:"source_account_id"`
	DestinationAccountID uint            `json:"destination_account_id"`
	Amount               decimal.Decimal `json:"amount"`
}
