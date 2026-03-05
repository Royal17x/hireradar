package utils

import (
	"database/sql"
	"github.com/pressly/goose/v3"
	"path/filepath"
	"runtime"
)

func RunMigrations(dsn string) error {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}

	_, filename, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(filename), "..", "..")
	migrationsDir := filepath.Join(projectRoot, "migrations")

	defer db.Close()
	return goose.Up(db, migrationsDir)
}
