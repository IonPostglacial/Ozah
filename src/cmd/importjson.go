package cmd

import (
	"flag"
	"fmt"
	"os"

	"nicolas.galipot.net/hazo/storage/dataset"
)

func ImportJson(args []string) error {
	fs := flag.NewFlagSet("importjson", flag.ExitOnError)

	var filePath, ds string
	fs.StringVar(&filePath, "file", "", "Path to the JSON file to import (required)")
	fs.StringVar(&ds, "dataset", "", "Name of the dataset to import into (required)")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: hazo importjson -file <path> -dataset <name>\n\n")
		fmt.Fprintf(fs.Output(), "Import data from a JSON file into a dataset.\n\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return err
	}

	if filePath == "" || ds == "" {
		fs.Usage()
		return fmt.Errorf("all flags are required: -file, -dataset")
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error importing file '%s': %w", filePath, err)
	}
	return dataset.ImportJson(data, dataset.Private(ds))
}
