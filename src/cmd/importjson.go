package cmd

import (
	"fmt"
	"os"

	"nicolas.galipot.net/hazo/storage"
)

func ImportJson(args []string) error {
	filePath := args[0]
	ds := args[1]
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error importing file '%s': %w", filePath, err)
	}
	return storage.ImportJson(data, storage.PrivateDataset(ds))
}
