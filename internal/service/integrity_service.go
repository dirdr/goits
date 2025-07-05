package service

import (
	"context"
	"fmt"

	"github.com/dirdr/goits/internal/domain"
	"github.com/dirdr/goits/internal/repository"
)

type integrityService struct {
	journalRepo repository.JournalRepository
}

func NewIntegrityService(journalRepo repository.JournalRepository) IntegrityService {
	return &integrityService{
		journalRepo: journalRepo,
	}
}

func (s *integrityService) VerifyDoubleBookkeeping(ctx context.Context) (*IntegrityResult, error) {
	totals, err := s.journalRepo.GetTotalsByEntryType(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get totals by entry type: %w", err)
	}

	totalDebits := totals[domain.Debit]
	totalCredits := totals[domain.Credit]
	difference := totalDebits.Sub(totalCredits)

	result := &IntegrityResult{
		IsValid:      difference.IsZero(),
		TotalDebits:  totalDebits,
		TotalCredits: totalCredits,
		Difference:   difference,
	}

	return result, nil
}
