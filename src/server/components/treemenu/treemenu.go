package treemenu

import (
	"context"
	_ "embed"
	"net/url"

	"nicolas.galipot.net/hazo/db"
	"nicolas.galipot.net/hazo/server/common"
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

type LangSet struct {
	common.BitSet
}

type Language uint64

const (
	Lang1 = Language(1 << iota)
	Lang2
	Lang3
)

type Lang struct {
	Name     string
	Url      string
	Selected bool
}

func LangSetFromString(s string) LangSet {
	return LangSet{common.BitSetFromString(s, common.BitSet(Lang1|Lang2|Lang3), common.BitSet(Lang1))}
}

func (lang LangSet) LangsFromNames(url *url.URL, names []string) []Lang {
	langs := make([]Lang, len(names))
	for i, name := range names {
		value := common.BitSet(1 << i)
		selected := lang.Contains(value)
		newLangs := lang.Toggle(value)
		query := url.Query()
		query.Del("menuLangs")
		query.Add("menuLangs", newLangs.String())
		newUrl := *url
		newUrl.RawQuery = query.Encode()
		langs[i] = Lang{
			Name:     name,
			Url:      newUrl.String(),
			Selected: selected,
		}
	}
	return langs
}

func LoadItemFromDb(ctx context.Context, queries *db.Queries, root string, langs []string, filter string) (*Item, error) {
	docs, err := queries.GetDocumentHierarchy(ctx, root, langs, filter)
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
