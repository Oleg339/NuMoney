package transaction

import (
	"context"
	"errors"
	"math"
	"numoney/internal/category"
	"time"
)

type Service struct {
	repository Repository
	catRepo    category.Repository
}

func NewService(rep Repository, catRepo category.Repository) Service {
	return Service{repository: rep, catRepo: catRepo}
}

var (
	ErrCategoryNotFound = errors.New("category not found")
)

func (s *Service) Create(ctx context.Context, amount float64, catID int, userID int64) (Transaction, error) {
	if _, err := s.catRepo.Find(ctx, int64(catID), userID); err != nil {
		return Transaction{}, category.ErrNotFound
	}

	transaction := &Transaction{
		Amount:     int(math.Round(amount * 100)),
		CategoryID: catID,
		UserID:     userID,
	}

	return *transaction, s.repository.Save(ctx, transaction)
}

func (s *Service) GetStatsForPeriod(ctx context.Context, userID int64, from, to time.Time) ([]Stat, error) {
	return s.repository.GetStatsForPeriodFromDB(ctx, userID, from, to)
}
