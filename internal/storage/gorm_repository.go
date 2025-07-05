package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/dirdr/goits/internal/domain"
	"github.com/shopspring/decimal"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormAccountRepository struct {
	db *gorm.DB
}

func NewGormAccountRepository(db *gorm.DB) *GormAccountRepository {
	return &GormAccountRepository{db: db}
}

func (r *GormAccountRepository) CreateAccount(ctx context.Context, tx *gorm.DB, account *domain.Account) error {
	gormAccount := GormAccount{
		ID:        account.ID,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
	}

	db := r.db
	if tx != nil {
		db = tx
	}

	result := db.WithContext(ctx).Create(&gormAccount)
	if result.Error != nil {
		return fmt.Errorf("failed to create account: %w", result.Error)
	}
	return nil
}

func (r *GormAccountRepository) GetAccountByID(ctx context.Context, tx *gorm.DB, accountID uint) (*domain.Account, error) {
	var gormAccount GormAccount

	db := r.db
	if tx != nil {
		db = tx
	}

	result := db.WithContext(ctx).First(&gormAccount, "id = ?", accountID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get account by ID: %w", result.Error)
	}

	return &domain.Account{
		ID:        gormAccount.ID,
		CreatedAt: gormAccount.CreatedAt,
		UpdatedAt: gormAccount.UpdatedAt,
	}, nil
}

func (r *GormAccountRepository) AccountExists(ctx context.Context, tx *gorm.DB, accountID uint) (bool, error) {
	var count int64

	db := r.db
	if tx != nil {
		db = tx
	}

	result := db.WithContext(ctx).Model(&GormAccount{}).Where("id = ?", accountID).Count(&count)
	if result.Error != nil {
		return false, fmt.Errorf("failed to check account existence: %w", result.Error)
	}

	return count > 0, nil
}

type GormAccountBalanceRepository struct {
	db *gorm.DB
}

func NewGormAccountBalanceRepository(db *gorm.DB) *GormAccountBalanceRepository {
	return &GormAccountBalanceRepository{db: db}
}

func (r *GormAccountBalanceRepository) GetAccountBalance(ctx context.Context, tx *gorm.DB, accountID uint) (*domain.AccountBalance, error) {
	var gormBalance GormAccountBalance

	db := r.db
	if tx != nil {
		db = tx
	}

	result := db.WithContext(ctx).First(&gormBalance, "account_id = ?", accountID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get account balance: %w", result.Error)
	}

	return &domain.AccountBalance{
		AccountID:   gormBalance.AccountID,
		Balance:     gormBalance.Balance,
		Version:     gormBalance.Version,
		LastEventID: gormBalance.LastEventID,
		UpdatedAt:   gormBalance.UpdatedAt,
	}, nil
}

func (r *GormAccountBalanceRepository) UpsertAccountBalance(ctx context.Context, tx *gorm.DB, balance *domain.AccountBalance) error {
	gormBalance := GormAccountBalance{
		AccountID:   balance.AccountID,
		Balance:     balance.Balance,
		Version:     balance.Version,
		LastEventID: balance.LastEventID,
		UpdatedAt:   balance.UpdatedAt,
	}

	db := r.db
	if tx != nil {
		db = tx
	}

	result := db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "account_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"balance", "version", "last_event_id", "updated_at"}),
		}).Create(&gormBalance)

	if result.Error != nil {
		return fmt.Errorf("failed to upsert account balance: %w", result.Error)
	}
	return nil
}

type GormTransferEventRepository struct {
	db *gorm.DB
}

func NewGormTransferEventRepository(db *gorm.DB) *GormTransferEventRepository {
	return &GormTransferEventRepository{db: db}
}

func (r *GormTransferEventRepository) SaveTransferEvent(ctx context.Context, tx *gorm.DB, event *domain.TransferEvent) error {
	gormEvent := GormTransferEvent{
		EventID:       event.EventID,
		TransferID:    event.TransferID,
		FromAccountID: event.FromAccountID,
		ToAccountID:   event.ToAccountID,
		Amount:        event.Amount,
		EventType:     event.EventType,
		CreatedAt:     event.CreatedAt,
	}

	db := r.db
	if tx != nil {
		db = tx
	}

	result := db.WithContext(ctx).Create(&gormEvent)
	if result.Error != nil {
		return fmt.Errorf("failed to save transfer event: %w", result.Error)
	}

	// Update the event ID in the domain object
	event.EventID = gormEvent.EventID
	return nil
}

type GormJournalRepository struct {
	db *gorm.DB
}

func NewGormJournalRepository(db *gorm.DB) *GormJournalRepository {
	return &GormJournalRepository{db: db}
}

func (r *GormJournalRepository) SaveJournalEntry(ctx context.Context, tx *gorm.DB, entry *domain.JournalEntry) error {
	gormEntry := GormJournalEntry{
		EntryID:       entry.EntryID,
		TransactionID: entry.TransactionID,
		AccountID:     entry.AccountID,
		Amount:        entry.Amount,
		Type:          entry.Type,
		SourceEventID: entry.SourceEventID,
		CreatedAt:     entry.CreatedAt,
	}

	db := r.db
	if tx != nil {
		db = tx
	}

	result := db.WithContext(ctx).Create(&gormEntry)
	if result.Error != nil {
		return fmt.Errorf("failed to save journal entry: %w", result.Error)
	}
	return nil
}

func (r *GormJournalRepository) GetTotalsByEntryType(ctx context.Context, tx *gorm.DB) (map[domain.EntryType]decimal.Decimal, error) {
	var results []struct {
		Type  domain.EntryType
		Total decimal.Decimal
	}

	db := r.db
	if tx != nil {
		db = tx
	}

	result := db.WithContext(ctx).
		Model(&GormJournalEntry{}).
		Select("type, SUM(amount) as total").
		Group("type").
		Find(&results)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get totals by entry type: %w", result.Error)
	}

	totals := make(map[domain.EntryType]decimal.Decimal)
	for _, r := range results {
		totals[r.Type] = r.Total
	}

	return totals, nil
}
