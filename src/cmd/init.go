package cmd

import "nicolas.galipot.net/hazo/db"

func Init(dbPath string) error {
	err := db.ExecSqlite(dbPath, db.Schema)
	if err != nil {
		return err
	}
	err = db.ExecSqlite(dbPath, db.Index)
	if err != nil {
		return err
	}
	return nil
}
