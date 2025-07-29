package sqlite

import (
	"database/sql"

	"github.com/surajNirala/student-api/internal/config"

	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}
	createQuery := `CREATE TABLE IF NOT EXISTS students (
						id INTEGER PRIMARY KEY AUTOINCREMENT,
						name TEXT,
						email TEXT,
						age INT,
						created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
						updated_at TIMESTAMP DEFAULT CURRECT_TIMESTAMP
					)`
	_, err = db.Exec(createQuery)
	if err != nil {
		return nil, err
	}
	return &Sqlite{Db: db}, nil
}
