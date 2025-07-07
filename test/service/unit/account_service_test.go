package unit

import (
	"context"
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

func TestAccountService_CreateAccount_Success(t *testing.T) {
	mockAccountRepo := &MockAccountRepository{}
	mockBalanceRepo := &MockAccountBalanceRepository{}
	tx := &gorm.DB{}

	mockAccountRepo.On("AccountExists", mock.Anything, tx, uint(1)).Return(false, nil)
	mockAccountRepo.On("CreateAccount", mock.Anything, tx, mock.AnythingOfType("*domain.Account")).Return(nil)
	mockBalanceRepo.On("UpsertAccountBalance", mock.Anything, tx, mock.AnythingOfType("*domain.AccountBalance")).Return(nil)

	svc := service.NewAccountService(mockAccountRepo, mockBalanceRepo)

	account, err := svc.CreateAccount(context.Background(), tx, 1, decimal.NewFromInt(100))

	require.NoError(t, err)
	assert.NotNil(t, account)
	assert.Equal(t, uint(1), account.ID)
	mockAccountRepo.AssertExpectations(t)
	mockBalanceRepo.AssertExpectations(t)
}

func TestAccountService_CreateAccount_NegativeBalance(t *testing.T) {
	mockAccountRepo := &MockAccountRepository{}
	mockBalanceRepo := &MockAccountBalanceRepository{}
	tx := &gorm.DB{}

	svc := service.NewAccountService(mockAccountRepo, mockBalanceRepo)

	account, err := svc.CreateAccount(context.Background(), tx, 1, decimal.NewFromInt(-10))

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "initial balance cannot be negative")
	assert.Nil(t, account)
}

func TestAccountService_GetAccountByID_Success(t *testing.T) {
	mockAccountRepo := &MockAccountRepository{}
	mockBalanceRepo := &MockAccountBalanceRepository{}

	expectedAccount := &domain.Account{
		ID:        1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockAccountRepo.On("GetAccountByID", mock.Anything, (*gorm.DB)(nil), uint(1)).Return(expectedAccount, nil)

	svc := service.NewAccountService(mockAccountRepo, mockBalanceRepo)

	account, err := svc.GetAccountByID(context.Background(), 1)

	require.NoError(t, err)
	assert.NotNil(t, account)
	assert.Equal(t, uint(1), account.ID)
	mockAccountRepo.AssertExpectations(t)
}

func TestAccountService_GetAccountBalance_Success(t *testing.T) {
	mockAccountRepo := &MockAccountRepository{}
	mockBalanceRepo := &MockAccountBalanceRepository{}

	expectedBalance := &domain.AccountBalance{
		AccountID:   1,
		Balance:     decimal.NewFromInt(100),
		Version:     1,
		LastEventID: 0,
		UpdatedAt:   time.Now(),
	}

	mockBalanceRepo.On("GetAccountBalance", mock.Anything, (*gorm.DB)(nil), uint(1)).Return(expectedBalance, nil)

	svc := service.NewAccountService(mockAccountRepo, mockBalanceRepo)

	balance, err := svc.GetAccountBalance(context.Background(), 1)

	require.NoError(t, err)
	assert.NotNil(t, balance)
	assert.Equal(t, uint(1), balance.AccountID)
	assert.Equal(t, decimal.NewFromInt(100), balance.Balance)
	mockBalanceRepo.AssertExpectations(t)
}
