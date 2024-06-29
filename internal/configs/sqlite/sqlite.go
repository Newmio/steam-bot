package sqlite

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

func OpenDb(dbName string) (*sqlx.DB, error) {
	return sqlx.Open("sqlite3", fmt.Sprintf("./database/%s.db", dbName))
}
