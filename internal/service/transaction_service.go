package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dirdr/goits/internal/domain"
	"github.com/dirdr/goits/internal/repository"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type transactionService struct {
	accountRepo        repository.AccountRepository
	accountBalanceRepo repository.AccountBalanceRepository
	transferEventRepo  repository.TransferEventRepository
	journalRepo        repository.JournalRepository
}

func NewTransactionService(
	accountRepo repository.AccountRepository,
	accountBalanceRepo repository.AccountBalanceRepository,
	transferEventRepo repository.TransferEventRepository,
	journalRepo repository.JournalRepository,
) TransactionService {
	return &transactionService{
		accountRepo:        accountRepo,
		accountBalanceRepo: accountBalanceRepo,
		transferEventRepo:  transferEventRepo,
		journalRepo:        journalRepo,
	}
}

func (s *transactionService) ProcessTransfer(ctx context.Context, tx *gorm.DB, sourceAccountID, destinationAccountID uint, amount decimal.Decimal) error {
	if amount.IsNegative() || amount.IsZero() {
		return errors.New("transfer amount must be positive")
	}

	if sourceAccountID == destinationAccountID {
		return errors.New("source and destination accounts cannot be the same")
	}

	// Check if accounts exist
	sourceExists, err := s.accountRepo.AccountExists(ctx, tx, sourceAccountID)
	if err != nil {
		return fmt.Errorf("failed to check source account: %w", err)
	}
	if !sourceExists {
		return errors.New("source account not found")
	}

	destinationExists, err := s.accountRepo.AccountExists(ctx, tx, destinationAccountID)
	if err != nil {
		return fmt.Errorf("failed to check destination account: %w", err)
	}
	if !destinationExists {
		return errors.New("destination account not found")
	}

	// Check source account balance
	sourceBalance, err := s.accountBalanceRepo.GetAccountBalance(ctx, tx, sourceAccountID)
	if err != nil {
		return fmt.Errorf("failed to get source account balance: %w", err)
	}
	if sourceBalance == nil {
		return errors.New("source account balance not found")
	}

	if sourceBalance.Balance.LessThan(amount) {
		return errors.New("insufficient balance in source account")
	}

	// Get destination account balance
	destinationBalance, err := s.accountBalanceRepo.GetAccountBalance(ctx, tx, destinationAccountID)
	if err != nil {
		return fmt.Errorf("failed to get destination account balance: %w", err)
	}
	if destinationBalance == nil {
		return errors.New("destination account balance not found")
	}

	now := time.Now()
	transferID := uuid.New().String()

	// 1. Save transfer event (source of truth)
	transferEvent := &domain.TransferEvent{
		TransferID:    transferID,
		FromAccountID: sourceAccountID,
		ToAccountID:   destinationAccountID,
		Amount:        amount,
		EventType:     "TransferProcessed",
		CreatedAt:     now,
	}
	err = s.transferEventRepo.SaveTransferEvent(ctx, tx, transferEvent)
	if err != nil {
		return fmt.Errorf("failed to save transfer event: %w", err)
	}

	// 2. Create journal entries (projections)
	debitEntry := &domain.JournalEntry{
		TransactionID: transferID,
		AccountID:     sourceAccountID,
		Amount:        amount,
		Type:          domain.Debit,
		SourceEventID: transferEvent.EventID,
		CreatedAt:     now,
	}
	err = s.journalRepo.SaveJournalEntry(ctx, tx, debitEntry)
	if err != nil {
		return fmt.Errorf("failed to save debit journal entry: %w", err)
	}

	creditEntry := &domain.JournalEntry{
		TransactionID: transferID,
		AccountID:     destinationAccountID,
		Amount:        amount,
		Type:          domain.Credit,
		SourceEventID: transferEvent.EventID,
		CreatedAt:     now,
	}
	err = s.journalRepo.SaveJournalEntry(ctx, tx, creditEntry)
	if err != nil {
		return fmt.Errorf("failed to save credit journal entry: %w", err)
	}

	// 3. Update account balances (projections)
	newSourceBalance := &domain.AccountBalance{
		AccountID:   sourceAccountID,
		Balance:     sourceBalance.Balance.Sub(amount),
		Version:     sourceBalance.Version + 1,
		LastEventID: transferEvent.EventID,
		UpdatedAt:   now,
	}
	err = s.accountBalanceRepo.UpsertAccountBalance(ctx, tx, newSourceBalance)
	if err != nil {
		return fmt.Errorf("failed to update source account balance: %w", err)
	}

	newDestinationBalance := &domain.AccountBalance{
		AccountID:   destinationAccountID,
		Balance:     destinationBalance.Balance.Add(amount),
		Version:     destinationBalance.Version + 1,
		LastEventID: transferEvent.EventID,
		UpdatedAt:   now,
	}
	err = s.accountBalanceRepo.UpsertAccountBalance(ctx, tx, newDestinationBalance)
	if err != nil {
		return fmt.Errorf("failed to update destination account balance: %w", err)
	}

	return nil
}
