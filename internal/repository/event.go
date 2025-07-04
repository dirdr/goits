package repository

import (
	"context"

	"github.com/dirdr/goits/internal/domain"
)

type EventRepository interface {
	Create(ctx context.Context, event *domain.Event) error
}
