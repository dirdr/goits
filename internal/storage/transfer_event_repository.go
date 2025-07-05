package storage

import (
	"context"
	"fmt"

	"github.com/dirdr/goits/internal/domain"
	"gorm.io/gorm"
)

type GormTransferEventRepository struct {
	db *gorm.DB
}

func NewGormTransferEventRepository(db *gorm.DB) *GormTransferEventRepository {
	return &GormTransferEventRepository{db: db}
}

func (repo *GormTransferEventRepository) SaveTransferEvent(ctx context.Context, tx *gorm.DB, event *domain.TransferEvent) error {
	gormEvent := GormTransferEvent{
		EventID:       event.EventID,
		TransferID:    event.TransferID,
		FromAccountID: event.FromAccountID,
		ToAccountID:   event.ToAccountID,
		Amount:        event.Amount,
		EventType:     event.EventType,
		CreatedAt:     event.CreatedAt,
	}

	db := repo.db
	if tx != nil {
		db = tx
	}

	result := db.WithContext(ctx).Create(&gormEvent)
	if result.Error != nil {
		return fmt.Errorf("failed to save transfer event: %w", result.Error)
	}

	event.EventID = gormEvent.EventID
	return nil
}
