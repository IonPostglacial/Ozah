package cmd

import "nicolas.galipot.net/hazo/storage"

func Init(args []string) error {
	dbPath := args[0]
	return storage.Create(storage.PrivateDataset(dbPath))
}
