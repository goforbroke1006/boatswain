package common

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func OpenDBConn() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./chat-blocks.db")
	if err != nil {
		return nil, err
	}

	schemaQuery, err := os.ReadFile("./db/schema.sql")
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(string(schemaQuery)); err != nil {
		return nil, err
	}

	return db, nil
}
