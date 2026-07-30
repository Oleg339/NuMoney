package category

import (
	"database/sql"
	"log"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) Save(c *Category) error {
	res, err := r.db.Exec("insert into categories (name, type, user_id) values (?,?,?)", c.Name, c.Type, c.UserId)

	if err != nil {
		log.Print(c)

		return err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return err
	}

	c.UserId = int(id)

	return nil
}

func (r *Repository) GetByUserId(id int) ([]Category, error) {
	var categories []Category

	rows, err := r.db.Query("select id, name, type from categories where user_id = ? order by type desc", id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var c Category
		rows.Scan(&c.ID, &c.Name, &c.Type)
		categories = append(categories, c)
	}

	return categories, nil
}

func (r *Repository) GetByUserIdAndType(id int, Type string) ([]Category, error) {
	var categories []Category
	rows, err := r.db.Query("select id, name, type from categories where user_id = ? and type = ?", id, Type)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var c Category
		rows.Scan(&c.ID, &c.Name, &c.Type)
		categories = append(categories, c)
	}

	return categories, nil
}
