package user

import (
	"database/sql"
	"errors"
	"log"
)

type Repository interface {
	Save(id int, user *User) error
	GetByTelegramId(id int64) (User, error)
}

type MySQLRepository struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) *MySQLRepository {
	return &MySQLRepository{db}
}

func (r MySQLRepository) Save(id int, user *User) error {
	result, err := r.db.Exec("insert into users (telegram_id) values (?)", id)

	if err != nil {
		return err
	}

	newId, err := result.LastInsertId()

	if err != nil {
		return err
	}

	*user = User{int(newId), int(id)}

	return nil
}

func (r *MySQLRepository) GetByTelegramId(id int64) (User, error) {
	var user User

	err := r.db.QueryRow("select id, telegram_id from users where telegram_id = ?", id).Scan(&user.ID, &user.TelegramId)

	switch {
	case err == nil:
		return user, nil

	case errors.Is(err, sql.ErrNoRows):
		if err := r.Save(int(id), &user); err != nil {
			return User{}, err
		}

		return user, nil

	default:
		return User{}, err
	}
}
