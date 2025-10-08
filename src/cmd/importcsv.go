package cmd

import "nicolas.galipot.net/hazo/storage/dataset"

func ImportCsv(args []string) error {
	csvPath := args[0]
	ds := args[1]
	return dataset.ImportCsv(csvPath, dataset.Private(ds))
}
