package taxons

import (
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"net/http"
	"slices"

	"nicolas.galipot.net/hazo/db"
	"nicolas.galipot.net/hazo/server/common"
	"nicolas.galipot.net/hazo/server/components"
	"nicolas.galipot.net/hazo/server/components/treemenu"
)

//go:embed taxons.html
var taxonPage string

type State struct {
	PageTitle         string
	DatasetName       string
	AvailableDatasets []db.Dataset
	MenuState         *treemenu.State
	SelectedTaxon     *FormData
}

func Handler(w http.ResponseWriter, r *http.Request, cc *common.Context) error {
	tmpl := components.NewTemplate()
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
		return err
	}
	if taxonId != "" {
		taxon, err = LoadFormDataFromDb(ctx, queries, taxonId)
		if err != nil {
			return err
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
		"documentUrl": func(taxon *treemenu.Item) string {
			return fmt.Sprintf("/ds/%s/taxons/%s", dbName, taxon.Id)
		},
	})
	tmpl, err = tmpl.Parse(taxonPage)
	if err != nil {
		return err
	}
	tmpl, err = tmpl.Parse(FormTemplate)
	if err != nil {
		return err
	}
	items, err := treemenu.LoadItemFromDb(ctx, queries)
	if err != nil {
		return err
	}
	datasets, err := db.ListDatasets()
	if err != nil {
		return err
	}
	w.Header().Add("Content-Type", "text/html")
	err = tmpl.Execute(w, State{
		PageTitle:         "Hazo",
		DatasetName:       dbName,
		AvailableDatasets: datasets,
		SelectedTaxon:     taxon,
		MenuState: &treemenu.State{
			Selected: taxon.Id,
			Root:     items,
		},
	})
	if err != nil {
		return err
	}
	return nil
}
