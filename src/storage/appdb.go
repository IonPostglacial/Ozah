package storage

import (
	"database/sql"
	"fmt"

	commonstorage "nicolas.galipot.net/hazo/storage/appdb"
)

func OpenAppDb() (*sql.DB, *commonstorage.Queries, error) {
	db, err := sql.Open("sqlite3", "file:common.db?_foreign_keys=on&cache=shared&mode=rwc")
	if err != nil {
		return nil, nil, fmt.Errorf("could not connect to the database: %w", err)
	}
	queries := commonstorage.New(db)
	return db, queries, nil
}
