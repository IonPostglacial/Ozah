package cmd

import (
	"flag"
	"fmt"

	"nicolas.galipot.net/hazo/storage/dataset"
)

func Init(args []string) error {
	fs := flag.NewFlagSet("init", flag.ExitOnError)

	var dbPath string
	fs.StringVar(&dbPath, "db", "", "Path where the new dataset database will be created (required)")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: hazo init -db <path>\n\n")
		fmt.Fprintf(fs.Output(), "Create a new dataset database at the specified path.\n\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return err
	}

	if dbPath == "" {
		fs.Usage()
		return fmt.Errorf("required flag -db not provided")
	}

	return dataset.Create(dataset.Private(dbPath))
}
