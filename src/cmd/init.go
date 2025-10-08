package cmd

import "nicolas.galipot.net/hazo/storage/dataset"

func Init(args []string) error {
	dbPath := args[0]
	return dataset.Create(dataset.Private(dbPath))
}
