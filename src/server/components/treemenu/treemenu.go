package treemenu

import (
	"context"
	_ "embed"
	"fmt"

	"nicolas.galipot.net/hazo/db/storage"
)

//go:embed entry.html
var EntryTemplate string

//go:embed treemenu.html
var Template string

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

func LoadItemFromDb(ctx context.Context, queries *storage.Queries) (*Item, error) {
	docs, err := queries.GetDocumentHierarchyTr2(ctx, storage.GetDocumentHierarchyTr2Params{
		Path: "t0", Lang1: "V", Lang2: "CN",
	})
	if err != nil {
		return nil, err
	}
	h := &Item{Id: "t0", Name: "<TOP>", FullPath: "t0"}
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
