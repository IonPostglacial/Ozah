package cmd

import "nicolas.galipot.net/hazo/db"

func Init(args []string) error {
	dbPath := args[0]
	return db.Init(db.PrivateDataset(dbPath))
}
