package cmd

import (
	"context"
	"fmt"

	"nicolas.galipot.net/hazo/db"
	"nicolas.galipot.net/hazo/db/storage"
)

func LsDoc(args []string) error {
	ctx := context.Background()
	dbPath := args[0]
	docPath := args[1]
	queries, err := db.Open(dbPath)
	if err != nil {
		return err
	}
	acanthaceae, err := queries.GetDocumentHierarchy(ctx, storage.GetDocumentHierarchyParams{Path: docPath})
	if err != nil {
		return err
	}
	for _, doc := range acanthaceae {
		fmt.Printf("%s\t%s\n", doc.Path, doc.Name)
	}
	return nil
}
