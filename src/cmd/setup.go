package cmd

import "nicolas.galipot.net/hazo/db"

func Setup(args []string) error {
	err := db.ExecSqlite("common.db", db.CommonSchema)
	if err != nil {
		return err
	}
	return nil
}
