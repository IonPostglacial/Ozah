package dataset

import (
	"database/sql"
	"fmt"

	storage "nicolas.galipot.net/hazo/storage/dsdb"
)

func ConnectDb(ds Private) (*sql.DB, error) {
	return sql.Open("sqlite3", "file:"+string(ds)+"?_foreign_keys=on&cache=shared&mode=rwc")
}

func OpenDb(ds Private) (*Queries, error) {
	db, err := ConnectDb(ds)
	if err != nil {
		return nil, fmt.Errorf("could not open the database: %w", err)
	}
	queries := &Queries{Queries: storage.New(db), db: db}
	return queries, nil
}

func FullPath(path string, ref string) string {
	if path == "" {
		return ref
	}
	return fmt.Sprintf("%s.%s", path, ref)
}
