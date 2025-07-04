package repository

import (
	"context"

	"github.com/dirdr/goits/internal/domain"
	"github.com/google/uuid"
)

type JournalRepository interface {
	CreateTransaction(ctx context.Context, transactionID uuid.UUID, entries ...*domain.JournalEntry) error
}
