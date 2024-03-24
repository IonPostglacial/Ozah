package cmd

import "nicolas.galipot.net/hazo/db"

func ImportCsv(args []string) error {
	csvPath := args[0]
	to := args[1]
	return db.ImportCsv(csvPath, to)
}
