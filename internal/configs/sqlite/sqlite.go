package sqlite

import (
	"github.com/jmoiron/sqlx"
)

func OpenDb() (*sqlx.DB, error) {
	return sqlx.Open("sqlite3", "./bot.db")
}
