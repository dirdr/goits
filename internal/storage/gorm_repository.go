package storage

import (
	"context"
	"errors"

	"github.com/dirdr/goits/internal/domain"
	"github.com/dirdr/goits/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type gormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) repository.JournalRepository {
	return &gormRepository{db: db}
}

func (r *gormRepository) CreateTransaction(ctx context.Context, transactionID uuid.UUID, entries ...*domain.JournalEntry) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, entry := range entries {
			if entry == nil {
				return errors.New("nil journal entry provided")
			}

			gormEntry := JournalEntry{
				TransactionID: transactionID,
				AccountID:     entry.AccountID,
				Amount:        entry.Amount,
			}

			if err := tx.Create(&gormEntry).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
