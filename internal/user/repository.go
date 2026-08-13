package user

import (
	"context"
	"database/sql"
)

type Repository interface {
	Save(ctx context.Context, id int64, user *User) error
	GetByTelegramID(ctx context.Context, id int64) (User, error)
}

type MySQLRepository struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) *MySQLRepository {
	return &MySQLRepository{db}
}

func (r *MySQLRepository) Save(ctx context.Context, id int64, user *User) error {
	result, err := r.db.ExecContext(ctx, "INSERT INTO users (telegram_id) VALUES (?)", id)

	if err != nil {
		return err
	}

	newID, err := result.LastInsertId()

	if err != nil {
		return err
	}

	*user = User{newID, id}

	return nil
}

func (r *MySQLRepository) GetByTelegramID(ctx context.Context, id int64) (User, error) {
	var user User

	err := r.db.QueryRowContext(ctx, "SELECT id, telegram_id FROM users WHERE telegram_id = ?", id).Scan(&user.ID, &user.TelegramID)

	if err != nil {
		return User{}, err
	}

	return user, nil
}
