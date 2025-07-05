package unit

import (
	"context"
	"testing"

	"github.com/dirdr/goits/internal/domain"
	"github.com/dirdr/goits/internal/service"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

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

func TestTransactionService_ProcessTransfer_Success(t *testing.T) {
	mockAccountRepo := &MockAccountRepository{}
	mockBalanceRepo := &MockAccountBalanceRepository{}
	mockEventRepo := &MockTransferEventRepository{}
	mockJournalRepo := &MockJournalRepository{}
	tx := &gorm.DB{}

	sourceBalance := &domain.AccountBalance{
		AccountID: 1,
		Balance:   decimal.NewFromInt(500),
		Version:   1,
	}

	destBalance := &domain.AccountBalance{
		AccountID: 2,
		Balance:   decimal.NewFromInt(200),
		Version:   1,
	}

	mockAccountRepo.On("AccountExists", mock.Anything, tx, uint(1)).Return(true, nil)
	mockAccountRepo.On("AccountExists", mock.Anything, tx, uint(2)).Return(true, nil)
	mockBalanceRepo.On("GetAccountBalance", mock.Anything, tx, uint(1)).Return(sourceBalance, nil)
	mockBalanceRepo.On("GetAccountBalance", mock.Anything, tx, uint(2)).Return(destBalance, nil)
	mockEventRepo.On("SaveTransferEvent", mock.Anything, tx, mock.AnythingOfType("*domain.TransferEvent")).Return(nil)
	mockJournalRepo.On("SaveJournalEntry", mock.Anything, tx, mock.AnythingOfType("*domain.JournalEntry")).Return(nil).Twice()
	mockBalanceRepo.On("UpsertAccountBalance", mock.Anything, tx, mock.AnythingOfType("*domain.AccountBalance")).Return(nil).Twice()

	svc := service.NewTransactionService(mockAccountRepo, mockBalanceRepo, mockEventRepo, mockJournalRepo)

	err := svc.ProcessTransfer(context.Background(), tx, 1, 2, decimal.NewFromInt(100))

	require.NoError(t, err)
	mockAccountRepo.AssertExpectations(t)
	mockBalanceRepo.AssertExpectations(t)
	mockEventRepo.AssertExpectations(t)
	mockJournalRepo.AssertExpectations(t)
}

func TestTransactionService_ProcessTransfer_NegativeAmount(t *testing.T) {
	mockAccountRepo := &MockAccountRepository{}
	mockBalanceRepo := &MockAccountBalanceRepository{}
	mockEventRepo := &MockTransferEventRepository{}
	mockJournalRepo := &MockJournalRepository{}
	tx := &gorm.DB{}

	svc := service.NewTransactionService(mockAccountRepo, mockBalanceRepo, mockEventRepo, mockJournalRepo)

	err := svc.ProcessTransfer(context.Background(), tx, 1, 2, decimal.NewFromInt(-50))

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "transfer amount must be positive")
}

func TestTransactionService_ProcessTransfer_ZeroAmount(t *testing.T) {
	mockAccountRepo := &MockAccountRepository{}
	mockBalanceRepo := &MockAccountBalanceRepository{}
	mockEventRepo := &MockTransferEventRepository{}
	mockJournalRepo := &MockJournalRepository{}
	tx := &gorm.DB{}

	svc := service.NewTransactionService(mockAccountRepo, mockBalanceRepo, mockEventRepo, mockJournalRepo)

	err := svc.ProcessTransfer(context.Background(), tx, 1, 2, decimal.Zero)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "transfer amount must be positive")
}

func TestTransactionService_ProcessTransfer_SameAccount(t *testing.T) {
	mockAccountRepo := &MockAccountRepository{}
	mockBalanceRepo := &MockAccountBalanceRepository{}
	mockEventRepo := &MockTransferEventRepository{}
	mockJournalRepo := &MockJournalRepository{}
	tx := &gorm.DB{}

	svc := service.NewTransactionService(mockAccountRepo, mockBalanceRepo, mockEventRepo, mockJournalRepo)

	err := svc.ProcessTransfer(context.Background(), tx, 1, 1, decimal.NewFromInt(100))

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "source and destination accounts cannot be the same")
}

func TestTransactionService_ProcessTransfer_InsufficientBalance(t *testing.T) {
	mockAccountRepo := &MockAccountRepository{}
	mockBalanceRepo := &MockAccountBalanceRepository{}
	mockEventRepo := &MockTransferEventRepository{}
	mockJournalRepo := &MockJournalRepository{}
	tx := &gorm.DB{}

	sourceBalance := &domain.AccountBalance{
		AccountID: 1,
		Balance:   decimal.NewFromInt(50),
		Version:   1,
	}

	mockAccountRepo.On("AccountExists", mock.Anything, tx, uint(1)).Return(true, nil)
	mockAccountRepo.On("AccountExists", mock.Anything, tx, uint(2)).Return(true, nil)
	mockBalanceRepo.On("GetAccountBalance", mock.Anything, tx, uint(1)).Return(sourceBalance, nil)

	svc := service.NewTransactionService(mockAccountRepo, mockBalanceRepo, mockEventRepo, mockJournalRepo)

	err := svc.ProcessTransfer(context.Background(), tx, 1, 2, decimal.NewFromInt(100))

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient balance")
	mockAccountRepo.AssertExpectations(t)
	mockBalanceRepo.AssertExpectations(t)
}

func TestTransactionService_ProcessTransfer_SourceAccountNotFound(t *testing.T) {
	mockAccountRepo := &MockAccountRepository{}
	mockBalanceRepo := &MockAccountBalanceRepository{}
	mockEventRepo := &MockTransferEventRepository{}
	mockJournalRepo := &MockJournalRepository{}
	tx := &gorm.DB{}

	mockAccountRepo.On("AccountExists", mock.Anything, tx, uint(1)).Return(false, nil)

	svc := service.NewTransactionService(mockAccountRepo, mockBalanceRepo, mockEventRepo, mockJournalRepo)

	err := svc.ProcessTransfer(context.Background(), tx, 1, 2, decimal.NewFromInt(100))

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "source account not found")
	mockAccountRepo.AssertExpectations(t)
}

func TestTransactionService_ProcessTransfer_DestinationAccountNotFound(t *testing.T) {
	mockAccountRepo := &MockAccountRepository{}
	mockBalanceRepo := &MockAccountBalanceRepository{}
	mockEventRepo := &MockTransferEventRepository{}
	mockJournalRepo := &MockJournalRepository{}
	tx := &gorm.DB{}

	mockAccountRepo.On("AccountExists", mock.Anything, tx, uint(1)).Return(true, nil)
	mockAccountRepo.On("AccountExists", mock.Anything, tx, uint(2)).Return(false, nil)

	svc := service.NewTransactionService(mockAccountRepo, mockBalanceRepo, mockEventRepo, mockJournalRepo)

	err := svc.ProcessTransfer(context.Background(), tx, 1, 2, decimal.NewFromInt(100))

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "destination account not found")
	mockAccountRepo.AssertExpectations(t)
}
