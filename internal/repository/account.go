package repository

import (
	"context"

	"github.com/dirdr/goits/internal/domain"
	"github.com/shopspring/decimal"
)

type AccountRepository interface {
	Create(ctx context.Context, account *domain.Account) error
	FindByID(ctx context.Context, id int64) (*domain.Account, error)
	GetBalance(ctx context.Context, id int64) (decimal.Decimal, error)
}
