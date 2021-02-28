package db

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func Migrate(db *sqlx.DB) error {
	driver, err := sqlite3.WithInstance(db.DB, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("with instance: %v", err)
	}

	_, path, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("get migration path")
	}

	migrationPath := "file://" + filepath.Join(filepath.Dir(path), "migrations")
	migrator, err := migrate.NewWithDatabaseInstance(migrationPath, "sqlite3", driver)

	if err != nil {
		return fmt.Errorf("migrate new: %v", err)
	}

	return migrator.Up()
}
