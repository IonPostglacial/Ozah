package treemenu

import (
	"context"
	_ "embed"

	"nicolas.galipot.net/hazo/storage"
)

type ViewModel struct {
	Selected     string
	Langs        []Lang
	ColumnsCount int
	Root         *Item
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

type Language uint64

type Lang struct {
	Name     string
	Ref      string
	Selected bool
}

func LoadItemFromDb(ctx context.Context, queries *storage.Queries, root string, langs []string, filter string) (*Item, error) {
	docs, err := queries.GetDocumentHierarchy(ctx, root, langs, filter)
	if err != nil {
		return nil, err
	}
	h := &Item{Id: root, Name: "<TOP>", FullPath: root}
	previous := h
	parent := h
	breadcrumb := []*Item{}
	for i := range docs {
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
		fullPath := storage.FullPath(doc.Path, doc.Ref)
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
