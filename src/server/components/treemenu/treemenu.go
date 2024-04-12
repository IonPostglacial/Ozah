package treemenu

import (
	"context"
	_ "embed"
	"strconv"

	"nicolas.galipot.net/hazo/db"
)

type State struct {
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

type LangSet uint8

const (
	Lang1 = LangSet(1 << iota)
	Lang2
	Lang3
)

type Lang struct {
	Name     string
	Value    LangSet
	Selected bool
}

func LangFromString(s string) LangSet {
	n, err := strconv.ParseUint(s, 10, 3)
	if err != nil {
		return LangSet(Lang1 | Lang2 | Lang3)
	}
	lang := LangSet(n & 0b00000111)
	if lang == 0 {
		lang |= Lang1
	}
	return lang
}

func (lang LangSet) Contains(other LangSet) bool {
	return lang&other != 0
}

func (lang LangSet) SelectedNames(names []string) []string {
	langNames := make([]string, 0, len(names))
	for i, name := range names {
		if lang.Contains(LangSet(1 << i)) {
			langNames = append(langNames, name)
		}
	}
	return langNames
}

func (lang LangSet) LangsFromNames(names []string) []Lang {
	langs := make([]Lang, len(names))
	for i, name := range names {
		value := LangSet(1 << i)
		langs[i] = Lang{
			Name:     name,
			Value:    value,
			Selected: lang.Contains(value),
		}
	}
	return langs
}

func LoadItemFromDb(ctx context.Context, queries *db.Queries, root string, langs []string, filter string) (*Item, error) {
	docs, err := queries.GetDocumentHierarchy(ctx, root, langs[1:], filter)
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
