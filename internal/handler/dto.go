package handler

import "github.com/shopspring/decimal"

type CreateAccountRequest struct {
	AccountID      string          `json:"account_id"`
	InitialBalance decimal.Decimal `json:"initial_balance"`
}

type GetAccountResponse struct {
	AccountID string          `json:"account_id"`
	Balance   decimal.Decimal `json:"balance"`
}

type CreateTransactionRequest struct {
	SourceAccountID      string          `json:"source_account_id"`
	DestinationAccountID string          `json:"destination_account_id"`
	Amount               decimal.Decimal `json:"amount"`
}
