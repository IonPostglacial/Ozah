package cmd

import (
	"context"
	"fmt"

	"nicolas.galipot.net/hazo/storage"
)

func LsDoc(args []string) error {
	ctx := context.Background()
	dbPath := args[0]
	docPath := args[1]
	queries, err := storage.OpenDsDb(storage.PrivateDataset(dbPath))
	if err != nil {
		return err
	}
	acanthaceae, err := queries.GetDocumentHierarchy(ctx, docPath, []string{}, "")
	if err != nil {
		return err
	}
	for _, doc := range acanthaceae {
		fmt.Printf("%s\t%s\n", doc.Path, doc.Name)
	}
	return nil
}
