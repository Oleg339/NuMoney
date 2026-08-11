package category

import (
	"context"
	"database/sql"
)

type Repository interface {
	GetByUserID(ctx context.Context, id int64) ([]Category, error)
	Find(ctx context.Context, id int64, userID int64) (Category, error)
	GetByUserIDAndType(ctx context.Context, id int64, catType string) ([]Category, error)
	Save(ctx context.Context, c *Category) error
}

type MySQLRepository struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) *MySQLRepository {
	return &MySQLRepository{db}
}

func (r *MySQLRepository) Save(ctx context.Context, c *Category) error {
	res, err := r.db.ExecContext(ctx, "INSERT INTO categories (name, type, user_id) VALUES (?,?,?)", c.Name, c.Type, c.UserID)

	if err != nil {
		return err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return err
	}

	c.ID = int(id)

	return nil
}

func (r *MySQLRepository) Find(ctx context.Context, id int64, userID int64) (Category, error) {
	var c Category

	row := r.db.QueryRowContext(ctx, "SELECT id, name, type FROM categories WHERE id = ? AND user_id = ?", id, userID)

	err := row.Scan(&c.ID, &c.Name, &c.Type)

	if err != nil {
		return c, err
	}

	return c, nil
}

func (r *MySQLRepository) GetByUserID(ctx context.Context, id int64) ([]Category, error) {
	var categories []Category

	rows, err := r.db.QueryContext(ctx, "SELECT id, name, type FROM categories WHERE user_id = ? ORDER BY type DESC", id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Type); err != nil {
			return nil, err
		}

		categories = append(categories, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *MySQLRepository) GetByUserIDAndType(ctx context.Context, id int64, catType string) ([]Category, error) {
	var categories []Category

	rows, err := r.db.QueryContext(ctx, "SELECT id, name, type FROM categories WHERE user_id = ? AND type = ?", id, catType)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Type); err != nil {
			return nil, err
		}

		categories = append(categories, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}
