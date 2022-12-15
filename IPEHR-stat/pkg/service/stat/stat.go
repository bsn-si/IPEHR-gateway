package stat

import (
	"fmt"
	"strconv"
	"time"

	"ipehr/stat/pkg/localDB"
)

type Service struct {
	db *localDB.DB
}

func NewService(db *localDB.DB) *Service {
	return &Service{db}
}

func (s *Service) GetPatientsCount(period string) (uint64, error) {
	start, end := resolvPeriod(period)

	count, err := s.db.StatPatientsCountGet(start, end)
	if err != nil {
		return 0, fmt.Errorf("db.StatPatientsCountGet error: %w", err)
	}

	return count, nil
}

func (s *Service) GetDocumentsCount(period string) (uint64, error) {
	start, end := resolvPeriod(period)

	count, err := s.db.StatDocumentsCountGet(start, end)
	if err != nil {
		return 0, fmt.Errorf("db.GetDocumentsCount error: %w", err)
	}

	return count, nil
}

func resolvPeriod(period string) (int64, int64) {
	if period == "" {
		return 0, 32503662000
	}

	i, _ := strconv.Atoi(period)
	if i < 202201 || i > 300000 {
		return 0, 0
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
