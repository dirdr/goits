package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/dirdr/goits/internal/domain"
	"gorm.io/gorm"
)

type GormAccountRepository struct {
	db *gorm.DB
}

func NewGormAccountRepository(db *gorm.DB) *GormAccountRepository {
	return &GormAccountRepository{db: db}
}

func (repo *GormAccountRepository) CreateAccount(ctx context.Context, tx *gorm.DB, account *domain.Account) error {
	gormAccount := GormAccount{
		ID:        account.ID,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
	}

	db := repo.db
	if tx != nil {
		db = tx
	}

	result := db.WithContext(ctx).Create(&gormAccount)
	if result.Error != nil {
		return fmt.Errorf("failed to create account: %w", result.Error)
	}
	return nil
}

func (repo *GormAccountRepository) GetAccountByID(ctx context.Context, tx *gorm.DB, accountID uint) (*domain.Account, error) {
	var gormAccount GormAccount

	db := repo.db
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

func (repo *GormAccountRepository) AccountExists(ctx context.Context, tx *gorm.DB, accountID uint) (bool, error) {
	var count int64

	db := repo.db
	if tx != nil {
		db = tx
	}

	result := db.WithContext(ctx).Model(&GormAccount{}).Where("id = ?", accountID).Count(&count)
	if result.Error != nil {
		return false, fmt.Errorf("failed to check account existence: %w", result.Error)
	}

	return count > 0, nil
}
