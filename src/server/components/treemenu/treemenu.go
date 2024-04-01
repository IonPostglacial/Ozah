package treemenu

import (
	"context"
	_ "embed"
	"fmt"

	"nicolas.galipot.net/hazo/db"
	"nicolas.galipot.net/hazo/db/storage"
)

type State struct {
	Selected string
	Root     *Item
}

type Item struct {
	Id       string
	Url      string
	FullPath string
	Order    int
	Name     string
	NameV    string
	NameCN   string
	Children []*Item
}

func LoadItemFromDb(ctx context.Context, queries *db.Queries, root string, langs [3]string) (*Item, error) {
	docs, err := queries.GetDocumentHierarchyTr2(ctx, storage.GetDocumentHierarchyTr2Params{
		Path: root, Lang1: langs[1], Lang2: langs[2],
	})
	if err != nil {
		return nil, err
	}
	h := &Item{Id: root, Name: "<TOP>", FullPath: root}
	previous := h
	parent := h
	breadcrumb := []*Item{}
	for i := 0; i < len(docs); i++ {
		doc := docs[i]
		switch {
		case doc.Path == previous.FullPath:
			parent = previous
			breadcrumb = append(breadcrumb, parent)
		case doc.Path != parent.FullPath:
			for doc.Path != parent.FullPath && len(breadcrumb) > 0 {
				breadcrumb = breadcrumb[:len(breadcrumb)-1]
				parent = breadcrumb[len(breadcrumb)-1]
			}
		}
		fullPath := fmt.Sprintf("%s.%s", doc.Path, doc.Ref)
		taxon := &Item{
			Id:       doc.Ref,
			FullPath: fullPath,
			Order:    int(doc.DocOrder),
			Name:     doc.Name, NameV: doc.NameTr1.String, NameCN: doc.NameTr2.String,
		}
		parent.Children = append(parent.Children, taxon)
		previous = taxon
	}
	return h, nil
}
