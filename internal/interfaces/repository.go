package interfaces

import "database/sql"

type Repository interface {
	*sql.DB
}
