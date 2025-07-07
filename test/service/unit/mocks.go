package unit

import (
	"context"

	"github.com/dirdr/goits/internal/domain"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockAccountRepository struct {
	mock.Mock
}

func (m *MockAccountRepository) CreateAccount(ctx context.Context, tx *gorm.DB, account *domain.Account) error {
	args := m.Called(ctx, tx, account)
	return args.Error(0)
}

func (m *MockAccountRepository) GetAccountByID(ctx context.Context, tx *gorm.DB, accountID uint) (*domain.Account, error) {
	args := m.Called(ctx, tx, accountID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Account), args.Error(1)
}

func (m *MockAccountRepository) AccountExists(ctx context.Context, tx *gorm.DB, accountID uint) (bool, error) {
	args := m.Called(ctx, tx, accountID)
	return args.Bool(0), args.Error(1)
}

type MockAccountBalanceRepository struct {
	mock.Mock
}

func (m *MockAccountBalanceRepository) GetAccountBalance(ctx context.Context, tx *gorm.DB, accountID uint) (*domain.AccountBalance, error) {
	args := m.Called(ctx, tx, accountID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.AccountBalance), args.Error(1)
}

func (m *MockAccountBalanceRepository) UpsertAccountBalance(ctx context.Context, tx *gorm.DB, balance *domain.AccountBalance) error {
	args := m.Called(ctx, tx, balance)
	return args.Error(0)
}

func (m *MockAccountBalanceRepository) UpdateAccountBalanceWithVersion(ctx context.Context, tx *gorm.DB, balance *domain.AccountBalance, expectedVersion int) error {
	args := m.Called(ctx, tx, balance, expectedVersion)
	return args.Error(0)
}

type MockTransferEventRepository struct {
	mock.Mock
}

func (m *MockTransferEventRepository) SaveTransferEvent(ctx context.Context, tx *gorm.DB, event *domain.TransferEvent) error {
	args := m.Called(ctx, tx, event)
	return args.Error(0)
}

type MockJournalRepository struct {
	mock.Mock
}

func (m *MockJournalRepository) SaveJournalEntry(ctx context.Context, tx *gorm.DB, entry *domain.JournalEntry) error {
	args := m.Called(ctx, tx, entry)
	return args.Error(0)
}

func (m *MockJournalRepository) GetTotalsByEntryType(ctx context.Context, tx *gorm.DB) (map[domain.EntryType]decimal.Decimal, error) {
	args := m.Called(ctx, tx)
	return args.Get(0).(map[domain.EntryType]decimal.Decimal), args.Error(1)
}
