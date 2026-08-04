package category

import (
	"database/sql"
	"log"
)

type Repository interface {
	GetByUserId(id int) ([]Category, error)
	GetByUserIdAndType(id int, Type string) ([]Category, error)
	Save(c *Category) error
}

type MySQLRepository struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) *MySQLRepository {
	return &MySQLRepository{db}
}

func (r *MySQLRepository) Save(c *Category) error {
	res, err := r.db.Exec("insert into categories (name, type, user_id) values (?,?,?)", c.Name, c.Type, c.UserId)

	if err != nil {
		log.Print(c)

		return err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return err
	}

	c.ID = int(id)

	return nil
}

func (r *MySQLRepository) GetByUserId(id int) ([]Category, error) {
	var categories []Category

	rows, err := r.db.Query("select id, name, type from categories where user_id = ? order by type desc", id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var c Category
		e := rows.Scan(&c.ID, &c.Name, &c.Type)

		if e != nil {
			return nil, e
		}

		categories = append(categories, c)
	}

	if e := rows.Err(); e != nil {
		return nil, e
	}

	return categories, nil
}

func (r *MySQLRepository) GetByUserIdAndType(id int, Type string) ([]Category, error) {
	var categories []Category

	rows, err := r.db.Query("select id, name, type from categories where user_id = ? and type = ?", id, Type)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var c Category
		e := rows.Scan(&c.ID, &c.Name, &c.Type)

		if e != nil {
			return nil, e
		}

		categories = append(categories, c)
	}

	if e := rows.Err(); e != nil {
		return nil, e
	}

	return categories, nil
}
