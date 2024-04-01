package db

import (
	"database/sql"
	"fmt"

	"nicolas.galipot.net/hazo/db/storage"
)

func Open(dbPath string) (*Queries, error) {
	db, err := sql.Open("sqlite3", "file:"+dbPath)
	if err != nil {
		return nil, fmt.Errorf("could not connect to the database: %w", err)
	}
	queries := &Queries{Queries: storage.New(db), db: db}
	return queries, nil
}
