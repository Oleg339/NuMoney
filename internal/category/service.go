package category

import (
	"context"
	"errors"
	"strings"
)

type Service struct {
	repo Repository
}

func NewService(rep Repository) Service {
	return Service{repo: rep}
}

const (
	TypeIncome  = "income"
	TypeExpense = "expense"
)

var (
	ErrEmptyName = errors.New("empty category name")
	ErrBadType   = errors.New("invalid category type")
	ErrNotFound   = errors.New("category not found")
)

func (s *Service) Create(ctx context.Context, userID int64, name, catType string) (Category, error) {
	name = strings.TrimSpace(name)

	if name == "" {
		return Category{}, ErrEmptyName
	}

	if err := s.validateType(catType); err != nil {
		return Category{}, err
	}

	c := Category{Name: name, UserID: userID, Type: catType}

	if err := s.repo.Save(ctx, &c); err != nil {
		return Category{}, err
	}

	return c, nil
}

func (s *Service) GetByUserID(ctx context.Context, userID int64) ([]Category, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *Service) GetByUserIDAndType(ctx context.Context, userID int64, catType string) ([]Category, error) {
	if err := s.validateType(catType); err != nil {
		return []Category{}, err
	}

	return s.repo.GetByUserIDAndType(ctx, userID, catType)
}

func (s *Service) Find(ctx context.Context, id int64, userID int64) (Category, error) {
	return s.repo.Find(ctx, id, userID)
}

func (s *Service) validateType(catType string) error {
	if catType != TypeExpense && catType != TypeIncome {
		return ErrBadType
	}

	return nil
}


