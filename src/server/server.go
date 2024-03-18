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
	"nicolas.galipot.net/hazo/server/components/treemenu"
	"nicolas.galipot.net/hazo/server/views/taxons"
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
	SelectedTaxon *taxons.FormData
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
	tmpl, err = tmpl.Parse(taxons.FormTemplate)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	queries, err := db.Open(fmt.Sprintf("%s.sq3", dbName))
	if err != nil {
		log.Fatal(err)
	}
	items, err := treemenu.LoadItemFromDb(ctx, queries)
	if err != nil {
		log.Fatal(err)
	}
	var taxon *taxons.FormData
	if taxonId != "" {
		taxon, err = taxons.LoadFormDataFromDb(ctx, queries, taxonId)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		taxon = &taxons.FormData{}
	}
	var buf strings.Builder
	err = tmpl.Execute(&buf, State{
		DatasetName:   dbName,
		SelectedTaxon: taxon,
		MenuRoot:      items,
	})
	if err != nil {
		log.Fatal(err)
	}
	http.ServeContent(w, r, "index.html", time.Now(), strings.NewReader(buf.String()))
}
