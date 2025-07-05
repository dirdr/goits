package storage

import (
	"context"
	"fmt"

	"github.com/dirdr/goits/internal/domain"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type GormJournalRepository struct {
	db *gorm.DB
}

func NewGormJournalRepository(db *gorm.DB) *GormJournalRepository {
	return &GormJournalRepository{db: db}
}

func (repo *GormJournalRepository) SaveJournalEntry(ctx context.Context, tx *gorm.DB, entry *domain.JournalEntry) error {
	gormEntry := GormJournalEntry{
		EntryID:       entry.EntryID,
		TransactionID: entry.TransactionID,
		AccountID:     entry.AccountID,
		Amount:        entry.Amount,
		Type:          entry.Type,
		SourceEventID: entry.SourceEventID,
		CreatedAt:     entry.CreatedAt,
	}

	db := repo.db
	if tx != nil {
		db = tx
	}

	result := db.WithContext(ctx).Create(&gormEntry)
	if result.Error != nil {
		return fmt.Errorf("failed to save journal entry: %w", result.Error)
	}
	return nil
}

func (repo *GormJournalRepository) GetTotalsByEntryType(ctx context.Context, tx *gorm.DB) (map[domain.EntryType]decimal.Decimal, error) {
	var results []struct {
		Type  domain.EntryType
		Total decimal.Decimal
	}

	db := repo.db
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
