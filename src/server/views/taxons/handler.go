package taxons

import (
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"slices"

	"nicolas.galipot.net/hazo/db"
	"nicolas.galipot.net/hazo/server/common"
	"nicolas.galipot.net/hazo/server/components/treemenu"
)

//go:embed taxons.html
var taxonPage string

type State struct {
	DatasetName       string
	AvailableDatasets []db.Dataset
	MenuState         *treemenu.State
	SelectedTaxon     *FormData
}

func Handler(w http.ResponseWriter, r *http.Request, cc *common.Context) {
	tmpl := template.New("taxons")
	dbName := r.PathValue("dsName")
	if dbName == "" {
		dbName = "plants"
	}
	taxonId := r.PathValue("id")
	var (
		taxon *FormData
		err   error
	)
	ctx := context.Background()
	queries, err := db.Open(fmt.Sprintf("%s.sq3", dbName))
	if err != nil {
		log.Fatal(err)
	}
	if taxonId != "" {
		taxon, err = LoadFormDataFromDb(ctx, queries, taxonId)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		taxon = &FormData{}
	}
	tmpl.Funcs(template.FuncMap{
		"selectedDoc": func() string {
			return taxon.Id
		},
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
	tmpl, err = tmpl.Parse(taxonPage)
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
	tmpl, err = tmpl.Parse(FormTemplate)
	if err != nil {
		log.Fatal(err)
	}
	items, err := treemenu.LoadItemFromDb(ctx, queries)
	if err != nil {
		log.Fatal(err)
	}
	datasets, err := db.ListDatasets()
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Add("Content-Type", "text/html")
	err = tmpl.Execute(w, State{
		DatasetName:       dbName,
		AvailableDatasets: datasets,
		SelectedTaxon:     taxon,
		MenuState: &treemenu.State{
			Selected: taxon.Id,
			Root:     items,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
