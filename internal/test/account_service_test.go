package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/dirdr/goits/internal/domain"
	"github.com/dirdr/goits/internal/service"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// Mock repositories
type MockAccountRepository struct {
	mock.Mock
}

func (m *MockAccountRepository) CreateAccount(ctx context.Context, tx *gorm.DB, account *domain.Account) error {
	args := m.Called(ctx, tx, account)
	return args.Error(0)
}

func (m *MockAccountRepository) GetAccountByID(ctx context.Context, tx *gorm.DB, accountID uint) (*domain.Account, error) {
	args := m.Called(ctx, tx, accountID)
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
	return args.Get(0).(*domain.AccountBalance), args.Error(1)
}

func (m *MockAccountBalanceRepository) UpsertAccountBalance(ctx context.Context, tx *gorm.DB, balance *domain.AccountBalance) error {
	args := m.Called(ctx, tx, balance)
	return args.Error(0)
}

type MockDB struct {
	mock.Mock
}

func (m *MockDB) WithContext(ctx context.Context) *gorm.DB {
	args := m.Called(ctx)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Transaction(fn func(*gorm.DB) error) error {
	args := m.Called(fn)
	return args.Error(0)
}

func TestAccountService_CreateAccount(t *testing.T) {
	tests := []struct {
		name           string
		accountID      uint
		initialBalance decimal.Decimal
		setupMocks     func(*MockAccountRepository, *MockAccountBalanceRepository, *MockDB)
		expectedError  string
	}{
		{
			name:           "successful account creation",
			accountID:      1,
			initialBalance: decimal.NewFromInt(100),
			setupMocks: func(accountRepo *MockAccountRepository, balanceRepo *MockAccountBalanceRepository, db *MockDB) {
				// Mock transaction behavior
				db.On("WithContext", mock.AnythingOfType("*context.emptyCtx")).Return(&gorm.DB{}).Once()
				db.On("Transaction", mock.AnythingOfType("func(*gorm.DB) error")).Return(nil).Run(func(args mock.Arguments) {
					fn := args.Get(0).(func(*gorm.DB) error)
					// Simulate transaction execution
					accountRepo.On("AccountExists", mock.Anything, mock.Anything, uint(1)).Return(false, nil).Once()
					accountRepo.On("CreateAccount", mock.Anything, mock.Anything, mock.AnythingOfType("*domain.Account")).Return(nil).Once()
					balanceRepo.On("UpsertAccountBalance", mock.Anything, mock.Anything, mock.AnythingOfType("*domain.AccountBalance")).Return(nil).Once()
					fn(&gorm.DB{})
				}).Once()
			},
			expectedError: "",
		},
		{
			name:           "negative initial balance",
			accountID:      1,
			initialBalance: decimal.NewFromInt(-10),
			setupMocks:     func(accountRepo *MockAccountRepository, balanceRepo *MockAccountBalanceRepository, db *MockDB) {},
			expectedError:  "initial balance cannot be negative",
		},
		{
			name:           "account already exists",
			accountID:      1,
			initialBalance: decimal.NewFromInt(100),
			setupMocks: func(accountRepo *MockAccountRepository, balanceRepo *MockAccountBalanceRepository, db *MockDB) {
				db.On("WithContext", mock.AnythingOfType("*context.emptyCtx")).Return(&gorm.DB{}).Once()
				db.On("Transaction", mock.AnythingOfType("func(*gorm.DB) error")).Return(errors.New("account with this ID already exists")).Run(func(args mock.Arguments) {
					fn := args.Get(0).(func(*gorm.DB) error)
					accountRepo.On("AccountExists", mock.Anything, mock.Anything, uint(1)).Return(true, nil).Once()
					fn(&gorm.DB{})
				}).Once()
			},
			expectedError: "account with this ID already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockAccountRepo := &MockAccountRepository{}
			mockBalanceRepo := &MockAccountBalanceRepository{}
			mockDB := &MockDB{}

			tt.setupMocks(mockAccountRepo, mockBalanceRepo, mockDB)

			// Create service
			svc := service.NewAccountService(mockAccountRepo, mockBalanceRepo, mockDB)

			// Execute
			account, err := svc.CreateAccount(context.Background(), tt.accountID, tt.initialBalance)

			// Assertions
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, account)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, account)
				assert.Equal(t, tt.accountID, account.ID)
				assert.WithinDuration(t, time.Now(), account.CreatedAt, time.Second)
			}

			// Verify mocks
			mockAccountRepo.AssertExpectations(t)
			mockBalanceRepo.AssertExpectations(t)
			mockDB.AssertExpectations(t)
		})
	}
}

func TestAccountService_GetAccountByID(t *testing.T) {
	tests := []struct {
		name          string
		accountID     uint
		setupMocks    func(*MockAccountRepository, *MockAccountBalanceRepository, *MockDB)
		expectedError string
	}{
		{
			name:      "successful account retrieval",
			accountID: 1,
			setupMocks: func(accountRepo *MockAccountRepository, balanceRepo *MockAccountBalanceRepository, db *MockDB) {
				expectedAccount := &domain.Account{
					ID:        1,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				accountRepo.On("GetAccountByID", mock.Anything, (*gorm.DB)(nil), uint(1)).Return(expectedAccount, nil).Once()
			},
			expectedError: "",
		},
		{
			name:      "account not found",
			accountID: 999,
			setupMocks: func(accountRepo *MockAccountRepository, balanceRepo *MockAccountBalanceRepository, db *MockDB) {
				accountRepo.On("GetAccountByID", mock.Anything, (*gorm.DB)(nil), uint(999)).Return((*domain.Account)(nil), errors.New("account not found")).Once()
			},
			expectedError: "account not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockAccountRepo := &MockAccountRepository{}
			mockBalanceRepo := &MockAccountBalanceRepository{}
			mockDB := &MockDB{}

			tt.setupMocks(mockAccountRepo, mockBalanceRepo, mockDB)

			// Create service
			svc := service.NewAccountService(mockAccountRepo, mockBalanceRepo, mockDB)

			// Execute
			account, err := svc.GetAccountByID(context.Background(), tt.accountID)

			// Assertions
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, account)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, account)
				assert.Equal(t, tt.accountID, account.ID)
			}

			// Verify mocks
			mockAccountRepo.AssertExpectations(t)
		})
	}
}

func TestAccountService_GetAccountBalance(t *testing.T) {
	accountID := uint(1)
	expectedBalance := &domain.AccountBalance{
		AccountID:   accountID,
		Balance:     decimal.NewFromInt(100),
		Version:     1,
		LastEventID: 0,
		UpdatedAt:   time.Now(),
	}

	mockAccountRepo := &MockAccountRepository{}
	mockBalanceRepo := &MockAccountBalanceRepository{}
	mockDB := &MockDB{}

	mockBalanceRepo.On("GetAccountBalance", mock.Anything, (*gorm.DB)(nil), accountID).Return(expectedBalance, nil).Once()

	svc := service.NewAccountService(mockAccountRepo, mockBalanceRepo, mockDB)

	balance, err := svc.GetAccountBalance(context.Background(), accountID)

	require.NoError(t, err)
	assert.NotNil(t, balance)
	assert.Equal(t, accountID, balance.AccountID)
	assert.Equal(t, decimal.NewFromInt(100), balance.Balance)

	mockBalanceRepo.AssertExpectations(t)
}