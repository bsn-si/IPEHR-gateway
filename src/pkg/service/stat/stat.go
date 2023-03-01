package stat

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

type PatientsRepository interface {
	StatPatientsCountGet(ctx context.Context, start, end int64) (uint64, error)
	StatDocumentsCountGet(ctx context.Context, start, end int64) (uint64, error)
}

type Service struct {
	repo PatientsRepository
}

func NewService(repo PatientsRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetPatientsCount(ctx context.Context, period string) (uint64, error) {
	start, end := resolvePeriod(period)

	count, err := s.repo.StatPatientsCountGet(ctx, start, end)
	if err != nil {
		return 0, fmt.Errorf("db.StatPatientsCountGet error: %w", err)
	}

	return count, nil
}

func (s *Service) GetDocumentsCount(ctx context.Context, period string) (uint64, error) {
	start, end := resolvePeriod(period)

	count, err := s.repo.StatDocumentsCountGet(ctx, start, end)
	if err != nil {
		return 0, fmt.Errorf("db.GetDocumentsCount error: %w", err)
	}

	return count, nil
}

func resolvePeriod(period string) (int64, int64) {
	if period == "" {
		return 0, 32503662000
	}

	i, _ := strconv.Atoi(period)
	if i < 202201 || i > 300000 {
		return 32503662000, 32503662000
	}

	year, _ := strconv.Atoi(period[0:4])
	month, _ := strconv.Atoi(period[4:6])

	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)

	if month < 12 {
		month++
	} else {
		month = 1
		year++
	}

	end := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)

	return start.Unix(), end.Unix()
}
