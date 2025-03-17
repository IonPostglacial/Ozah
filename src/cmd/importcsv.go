package cmd

import "nicolas.galipot.net/hazo/storage"

func ImportCsv(args []string) error {
	csvPath := args[0]
	ds := args[1]
	return storage.ImportCsv(csvPath, storage.PrivateDataset(ds))
}
