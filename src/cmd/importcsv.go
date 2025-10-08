package cmd

import (
	"flag"
	"fmt"

	"nicolas.galipot.net/hazo/storage/dataset"
)

func ImportCsv(args []string) error {
	fs := flag.NewFlagSet("importcsv", flag.ExitOnError)

	var csvPath, ds string
	fs.StringVar(&csvPath, "csv", "", "Path to the CSV file to import (required)")
	fs.StringVar(&ds, "dataset", "", "Name of the dataset to import into (required)")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: hazo importcsv -csv <file> -dataset <name>\n\n")
		fmt.Fprintf(fs.Output(), "Import data from a CSV file into a dataset.\n\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return err
	}

	if csvPath == "" || ds == "" {
		fs.Usage()
		return fmt.Errorf("all flags are required: -csv, -dataset")
	}

	return dataset.ImportCsv(csvPath, dataset.Private(ds))
}
