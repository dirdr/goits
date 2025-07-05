package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/dirdr/goits/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormAccountBalanceRepository struct {
	db *gorm.DB
}

func NewGormAccountBalanceRepository(db *gorm.DB) *GormAccountBalanceRepository {
	return &GormAccountBalanceRepository{db: db}
}

func (repo *GormAccountBalanceRepository) GetAccountBalance(ctx context.Context, tx *gorm.DB, accountID uint) (*domain.AccountBalance, error) {
	var gormBalance GormAccountBalance

	db := repo.db
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

func (repo *GormAccountBalanceRepository) UpsertAccountBalance(ctx context.Context, tx *gorm.DB, balance *domain.AccountBalance) error {
	gormBalance := GormAccountBalance{
		AccountID:   balance.AccountID,
		Balance:     balance.Balance,
		Version:     balance.Version,
		LastEventID: balance.LastEventID,
		UpdatedAt:   balance.UpdatedAt,
	}

	db := repo.db
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
