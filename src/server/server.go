package server

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"

	"embed"

	"nicolas.galipot.net/hazo/db"
	"nicolas.galipot.net/hazo/db/storage"
	"nicolas.galipot.net/hazo/server/components/treemenu"
	"nicolas.galipot.net/hazo/server/views"
)

//go:embed assets
var assets embed.FS

//go:embed index.html
var indexPage string

func New() *http.ServeMux {
	s := http.NewServeMux()
	s.HandleFunc("/ds/{dsName}/taxons/{id}", indexHandler)
	s.HandleFunc("/", indexHandler)
	s.Handle("/assets/", http.FileServer(http.FS(assets)))
	return s
}

type State struct {
	DatasetName   string
	MenuRoot      *treemenu.Item
	SelectedTaxon *views.TaxonFormData
}

func taxonHierarchyFromDb(ctx context.Context, queries *storage.Queries) (*treemenu.Item, error) {
	docs, err := queries.GetDocumentHierarchyTr2(ctx, storage.GetDocumentHierarchyTr2Params{
		Path: "t0", Lang1: "V", Lang2: "CN",
	})
	if err != nil {
		return nil, err
	}
	h := &treemenu.Item{Id: "t0", Name: "<TOP>", FullPath: "t0"}
	previous := h
	parent := h
	breadcrumb := []*treemenu.Item{}
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
		taxon := &treemenu.Item{
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

func taxonFormDataFromDb(ctx context.Context, queries *storage.Queries, id string) (*views.TaxonFormData, error) {
	data, err := queries.GetTaxonInfo(ctx, id)
	if err != nil {
		return nil, err
	}
	return &views.TaxonFormData{
		Name:        data.Name,
		NameV:       data.NameV.String,
		NameCN:      data.NameCn.String,
		Description: data.Details.String,
		Author:      data.Author,
		Website:     data.Website.String,
	}, nil
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.New("index")
	dbName := r.PathValue("dsName")
	if dbName == "" {
		dbName = "plants"
	}
	taxonId := r.PathValue("id")
	tmpl.Funcs(template.FuncMap{
		"sortDocs": func(items []*treemenu.Item) []*treemenu.Item {
			slices.SortFunc(items, func(i, o *treemenu.Item) int {
				return i.Order - o.Order
			})
			return items
		},
		"taxonUrl": func(taxon *treemenu.Item) string {
			return fmt.Sprintf("/ds/%s/taxons/%s", dbName, taxon.Id)
		},
	})
	tmpl, err := tmpl.Parse(indexPage)
	if err != nil {
		log.Fatal(err)
	}
	tmpl, err = tmpl.Parse(treemenu.Template)
	if err != nil {
		log.Fatal(err)
	}
	tmpl, err = tmpl.Parse(treemenu.EntryTemplate)
	if err != nil {
		log.Fatal(err)
	}
	tmpl, err = tmpl.Parse(views.TaxonFormTemplate)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	queries, err := db.Open(fmt.Sprintf("%s.sq3", dbName))
	if err != nil {
		log.Fatal(err)
	}
	taxons, err := taxonHierarchyFromDb(ctx, queries)
	if err != nil {
		log.Fatal(err)
	}
	var taxon *views.TaxonFormData
	if taxonId != "" {
		taxon, err = taxonFormDataFromDb(ctx, queries, taxonId)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		taxon = &views.TaxonFormData{}
	}
	var buf strings.Builder
	err = tmpl.Execute(&buf, State{
		DatasetName:   dbName,
		SelectedTaxon: taxon,
		MenuRoot:      taxons,
	})
	if err != nil {
		log.Fatal(err)
	}
	http.ServeContent(w, r, "index.html", time.Now(), strings.NewReader(buf.String()))
}
