package user

import (
	"database/sql"
	"log"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return Repository{db}
}

func (r Repository) save(id int, user *User) error {
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

func (r *Repository) GetByTelegramId(id int64) (User, error) {
	var user User
	log.Print("ID: ", id)
	error := r.db.QueryRow("select id, telegram_id from users where telegram_id = ?", id).Scan(&user.ID, &user.TelegramId)

	if error == sql.ErrNoRows {
		r.save(int(id), &user)
	}

	log.Print("User: ", user)

	return user, nil
}
