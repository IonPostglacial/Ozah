package db

import (
	"database/sql"
	"fmt"

	"nicolas.galipot.net/hazo/db/commonstorage"
)

func OpenCommon() (*sql.DB, *commonstorage.Queries, error) {
	db, err := sql.Open("sqlite3", "file:common.db?_foreign_keys=on")
	if err != nil {
		return nil, nil, fmt.Errorf("could not connect to the database: %w", err)
	}
	queries := commonstorage.New(db)
	return db, queries, nil
}
