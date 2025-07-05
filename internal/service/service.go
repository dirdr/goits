package service

import (
	"context"

	"github.com/dirdr/goits/internal/domain"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type AccountService interface {
	CreateAccount(ctx context.Context, tx *gorm.DB, accountID uint, initialBalance decimal.Decimal) (*domain.Account, error)
	GetAccountByID(ctx context.Context, accountID uint) (*domain.Account, error)
	GetAccountBalance(ctx context.Context, accountID uint) (*domain.AccountBalance, error)
}

type TransactionService interface {
	ProcessTransfer(ctx context.Context, tx *gorm.DB, sourceAccountID, destinationAccountID uint, amount decimal.Decimal) error
}

type IntegrityService interface {
	VerifyDoubleBookkeeping(ctx context.Context) (*IntegrityResult, error)
}

type IntegrityResult struct {
	IsValid      bool            `json:"is_valid"`
	TotalDebits  decimal.Decimal `json:"total_debits"`
	TotalCredits decimal.Decimal `json:"total_credits"`
	Difference   decimal.Decimal `json:"difference"`
}
