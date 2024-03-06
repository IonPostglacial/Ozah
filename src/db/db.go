package db

import (
	"database/sql"
	"fmt"

	"nicolas.galipot.net/hazo/db/storage"
)

func Open(dbPath string) (*storage.Queries, error) {
	db, err := sql.Open("sqlite3", "file:"+dbPath)
	if err != nil {
		return nil, fmt.Errorf("could not connect to the database: %w", err)
	}
	queries := storage.New(db)
	return queries, nil
}
