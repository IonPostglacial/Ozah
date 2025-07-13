package cmd

import (
	"fmt"
	"os"

	"nicolas.galipot.net/hazo/storage"
)

func ExportJson(args []string) error {
	filePath := args[0]
	dsName := args[1]
	ds := storage.PrivateDataset(dsName)
	queries, err := storage.OpenDsDb(ds)
	if err != nil {
		return fmt.Errorf("could not open dataset database for '%s': %w", dsName, err)
	}
	outputFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("could not create output file '%s': %w", filePath, err)
	}
	defer outputFile.Close()
	err = storage.ExportJson(dsName, queries, outputFile)
	if err != nil {
		return fmt.Errorf("could not export dataset '%s' to JSON: %w", dsName, err)
	}
	return nil
}
