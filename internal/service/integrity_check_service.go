package service

import (
	"context"
	"fmt"

	"github.com/dirdr/goits/internal/domain"
	"github.com/dirdr/goits/internal/repository"
)

type integrityCheckService struct {
	journalRepo repository.JournalRepository
}

func NewIntegrityCheckService(journalRepo repository.JournalRepository) IntegrityCheckService {
	return &integrityCheckService{
		journalRepo: journalRepo,
	}
}

func (s *integrityCheckService) VerifyDoubleBookkeeping(ctx context.Context) (*IntegrityCheckResult, error) {
	totals, err := s.journalRepo.GetTotalsByEntryType(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get totals by entry type: %w", err)
	}

	totalDebits := totals[domain.Debit]
	totalCredits := totals[domain.Credit]
	difference := totalDebits.Sub(totalCredits)

	result := &IntegrityCheckResult{
		IsValid:      difference.IsZero(),
		TotalDebits:  totalDebits,
		TotalCredits: totalCredits,
		Difference:   difference,
	}

	return result, nil
}
