package cmd

import (
	"flag"
	"fmt"
	"os"

	"nicolas.galipot.net/hazo/storage/dataset"
)

func ExportJson(args []string) error {
	fs := flag.NewFlagSet("exportjson", flag.ExitOnError)

	var filePath, dsName string
	fs.StringVar(&filePath, "output", "", "Path where the JSON file will be created (required)")
	fs.StringVar(&dsName, "dataset", "", "Name of the dataset to export (required)")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: hazo exportjson -output <file> -dataset <name>\n\n")
		fmt.Fprintf(fs.Output(), "Export a dataset to JSON format.\n\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return err
	}

	if filePath == "" || dsName == "" {
		fs.Usage()
		return fmt.Errorf("all flags are required: -output, -dataset")
	}
	ds := dataset.Private(dsName)
	queries, err := dataset.OpenDb(ds)
	if err != nil {
		return fmt.Errorf("could not open dataset database for '%s': %w", dsName, err)
	}
	outputFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("could not create output file '%s': %w", filePath, err)
	}
	defer outputFile.Close()
	err = dataset.ExportJson(dsName, queries, outputFile)
	if err != nil {
		return fmt.Errorf("could not export dataset '%s' to JSON: %w", dsName, err)
	}
	return nil
}
