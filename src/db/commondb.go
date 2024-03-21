package db

import (
	"database/sql"
	"fmt"

	"nicolas.galipot.net/hazo/db/commonstorage"
)

func OpenCommon() (*commonstorage.Queries, error) {
	db, err := sql.Open("sqlite3", "file:common.db")
	if err != nil {
		return nil, fmt.Errorf("could not connect to the database: %w", err)
	}
	queries := commonstorage.New(db)
	return queries, nil
}
