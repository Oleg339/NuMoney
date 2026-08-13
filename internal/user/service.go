package user

import (
	"context"
	"database/sql"
	"errors"
)

type Service struct {
	repository Repository
}

func NewService(rep Repository) Service {
	return Service{repository: rep}
}

func (s *Service) EnsureRegistered(ctx context.Context, tgID int64) (User, error) {
	user, err := s.repository.GetByTelegramID(ctx, tgID)

	switch {
	case err == nil:
		return user, nil

	case errors.Is(err, sql.ErrNoRows):
		err := s.repository.Save(ctx, tgID, &user)

		return user, err
	}

	return User{}, err
}
