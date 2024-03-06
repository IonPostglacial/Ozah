package cmd

import (
	"context"
	"fmt"

	"nicolas.galipot.net/hazo/db"
	"nicolas.galipot.net/hazo/db/storage"
)

func LsDoc(ctx context.Context, dbPath string, docPath string) error {
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
