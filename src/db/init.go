package db

import (
	_ "embed"
	"fmt"
)

//go:embed schema.sql
var schema string

//go:embed index.sql
var index string

func Init(dbPath string) error {
	err := ExecSqlite(dbPath, fmt.Sprintf("%s\n.exit\n", schema))
	if err != nil {
		return err
	}
	err = ExecSqlite(dbPath, fmt.Sprintf("%s\n.exit\n", index))
	if err != nil {
		return err
	}
	return nil
}
