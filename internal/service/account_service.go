package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dirdr/goits/internal/domain"
	"github.com/dirdr/goits/internal/repository"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type accountService struct {
	accountRepo        repository.AccountRepository
	accountBalanceRepo repository.AccountBalanceRepository
}

func NewAccountService(accountRepo repository.AccountRepository, accountBalanceRepo repository.AccountBalanceRepository) AccountService {
	return &accountService{
		accountRepo:        accountRepo,
		accountBalanceRepo: accountBalanceRepo,
	}
}

func (s *accountService) CreateAccount(ctx context.Context, tx *gorm.DB, accountID uint, initialBalance decimal.Decimal) (*domain.Account, error) {
	if initialBalance.IsNegative() {
		return nil, errors.New("initial balance cannot be negative")
	}

	exists, err := s.accountRepo.AccountExists(ctx, tx, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to check for existing account: %w", err)
	}
	if exists {
		return nil, errors.New("account with this ID already exists")
	}

	account := &domain.Account{
		ID:        accountID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = s.accountRepo.CreateAccount(ctx, tx, account)
	if err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	balance := &domain.AccountBalance{
		AccountID:   accountID,
		Balance:     initialBalance,
		Version:     1,
		LastEventID: 0,
		UpdatedAt:   time.Now(),
	}

	err = s.accountBalanceRepo.UpsertAccountBalance(ctx, tx, balance)
	if err != nil {
		return nil, fmt.Errorf("failed to create initial balance: %w", err)
	}

	return account, nil
}

func (s *accountService) GetAccountByID(ctx context.Context, accountID uint) (*domain.Account, error) {
	account, err := s.accountRepo.GetAccountByID(ctx, nil, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}
	return account, nil
}

func (s *accountService) GetAccountBalance(ctx context.Context, accountID uint) (*domain.AccountBalance, error) {
	balance, err := s.accountBalanceRepo.GetAccountBalance(ctx, nil, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get account balance: %w", err)
	}
	return balance, nil
}