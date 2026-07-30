package transaction

import (
	"database/sql"
	"log"
	"time"
)

type Repository struct {
	db *sql.DB
}

type Stat struct {
	Name   string
	Type   string
	Amount int
}

func NewRepository(db *sql.DB) Repository {
	return Repository{db}
}

func (r *Repository) Save(t *Transaction) error {
	res, err := r.db.Exec("insert into transactions (category_id, user_id, amount, comment) values (?,?,?,?)", t.CategoryId, t.UserId, t.Amount, t.Comment)

	if err != nil {
		log.Fatal(err)
	}

	id, err := res.LastInsertId()

	t.ID = int(id)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetStatsForPeriodFromDb(userId int, from time.Time, to time.Time) ([]Stat, error) {
	rows, err := r.db.Query(`
          SELECT c.name, c.type, SUM(t.amount)
          FROM transactions t
          JOIN categories c ON c.id = t.category_id
          WHERE t.user_id = ? AND t.created_at BETWEEN ? AND ? 
          GROUP BY c.id order by type asc
      `, userId, from, to)

	if err != nil {
		log.Fatal(err)
	}

	var stats []Stat

	for rows.Next() {
		var stat Stat
		err := rows.Scan(&stat.Name, &stat.Type, &stat.Amount)

		if err != nil {
			log.Println("Scan error:", err)
		}
		stats = append(stats, stat)
	}

	if e := rows.Err(); e != nil {
		return nil, e
	}

	return stats, nil
}
