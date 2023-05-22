package common

import (
	"database/sql"
	"os"
)

func OpenDBConn(sourcePath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", sourcePath)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func ApplyMigrationFile(db *sql.DB, path string) error {
	query, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if _, err := db.Exec(string(query)); err != nil {
		return err
	}
	return nil
}
