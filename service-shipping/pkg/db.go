package pkg

import (
	"database/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func InitDB() (*sql.DB, error) {
	env, err := LoadEnv()
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("pgx", env["DB_DSN"])
	if err != nil {
		return nil, err
	}
	return db, nil
}
