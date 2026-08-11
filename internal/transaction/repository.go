package transaction

import (
	"context"
	"database/sql"
	"time"
)

type Repository interface {
	Save(ctx context.Context, t *Transaction) error
	GetStatsForPeriodFromDB(ctx context.Context, userID int64, from time.Time, to time.Time) ([]Stat, error)
}

type MySQLRepository struct {
	db *sql.DB
}

type Stat struct {
	Name   string
	Type   string
	Amount int
}

func NewMySQLRepository(db *sql.DB) *MySQLRepository {
	return &MySQLRepository{db}
}

func (r *MySQLRepository) Save(ctx context.Context, t *Transaction) error {
	res, err := r.db.ExecContext(ctx, "INSERT INTO transactions (category_id, user_id, amount, comment) VALUES (?,?,?,?)", t.CategoryID, t.UserID, t.Amount, t.Comment)

	if err != nil {
		return err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return err
	}

	t.ID = int(id)

	return nil
}

func (r *MySQLRepository) GetStatsForPeriodFromDB(ctx context.Context, userID int64, from time.Time, to time.Time) ([]Stat, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT c.name, c.type, SUM(t.amount)
		FROM transactions t
		JOIN categories c ON c.id = t.category_id
		WHERE t.user_id = ? AND t.created_at BETWEEN ? AND ?
		GROUP BY c.id ORDER BY type ASC
	`, userID, from, to)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var stats []Stat

	for rows.Next() {
		var stat Stat
		if err := rows.Scan(&stat.Name, &stat.Type, &stat.Amount); err != nil {
			return nil, err
		}

		stats = append(stats, stat)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stats, nil
}
