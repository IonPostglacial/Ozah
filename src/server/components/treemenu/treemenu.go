package treemenu

import (
	"context"
	_ "embed"

	"nicolas.galipot.net/hazo/db"
)

type State struct {
	Selected string
	Langs    []string
	Root     *Item
}

type Item struct {
	Id       string
	Url      string
	FullPath string
	Order    int64
	Name     string
	NameTr   []string
	Children []*Item
}

func LoadItemFromDb(ctx context.Context, queries *db.Queries, root string, langs [3]string) (*Item, error) {
	docs, err := queries.GetDocumentHierarchy(ctx, root, langs[1:], "")
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
				if len(breadcrumb) > 0 {
					parent = breadcrumb[len(breadcrumb)-1]
				}
			}
		}
		fullPath := db.FullPath(doc.Path, doc.Ref)
		nameTr := make([]string, len(doc.NameTr))
		for i, name := range doc.NameTr {
			nameTr[i] = name.String
		}
		taxon := &Item{
			Id:       doc.Ref,
			FullPath: fullPath,
			Order:    doc.DocOrder,
			Name:     doc.Name,
			NameTr:   nameTr,
		}
		parent.Children = append(parent.Children, taxon)
		previous = taxon
	}
	return h, nil
}
