package db

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
)

func NewDB(dir, name string) (*sqlx.DB, error) {
	file, err := os.Create(filepath.Join(dir, name))
	if err != nil {
		return nil, fmt.Errorf("create sqlite file: %v", err)
	}
	file.Close()

	db, err := sqlx.Open("sqlite3", file.Name())
	if err != nil {
		return nil, fmt.Errorf("open sqlite file: %v", err)
	}

	return db, err
}
