package cmd

import (
	"context"
	"flag"
	"fmt"

	"nicolas.galipot.net/hazo/storage/dataset"
)

func LsDoc(args []string) error {
	fs := flag.NewFlagSet("lsdoc", flag.ExitOnError)

	var dbPath, docPath string
	fs.StringVar(&dbPath, "db", "", "Name of the dataset database (required)")
	fs.StringVar(&docPath, "path", "", "Path to the document hierarchy to list (required)")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: hazo lsdoc -db <database> -path <document-path>\n\n")
		fmt.Fprintf(fs.Output(), "List documents in a dataset.\n\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return err
	}

	if dbPath == "" || docPath == "" {
		fs.Usage()
		return fmt.Errorf("all flags are required: -db, -path")
	}

	ctx := context.Background()
	queries, err := dataset.OpenDb(dataset.Private(dbPath))
	if err != nil {
		return err
	}
	hierarchy, err := queries.GetDocumentHierarchy(ctx, docPath, []string{}, "")
	if err != nil {
		return err
	}
	for _, doc := range hierarchy {
		fmt.Printf("%s\t%s\n", doc.Path, doc.Name)
	}
	return nil
}
