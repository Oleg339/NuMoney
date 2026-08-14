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
	ErrUnknownPeriod    = errors.New("unknown period")
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

func resolvePeriod(kind string) (from, to time.Time, err error) {
      now := time.Now()

      switch kind {
      case "today_stats":
			from = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			to = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())

      case "week_stats":
			weekday := int(now.Weekday())
			if weekday == 0 {
					weekday = 7
			}
			from = time.Date(now.Year(), now.Month(), now.Day()-weekday+1, 0, 0, 0, 0, now.Location())
			to = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())

      case "month_stats":
			from = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
			to = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())

      case "all_stats":
			from = time.Date(2000, 1, 1, 0, 0, 0, 0, now.Location())
			to = now

      default:
			err = ErrUnknownPeriod
      }

      return
}
