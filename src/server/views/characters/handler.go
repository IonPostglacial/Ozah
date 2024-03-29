package characters

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
	"nicolas.galipot.net/hazo/server/components/popover"
	"nicolas.galipot.net/hazo/server/components/treemenu"
	"nicolas.galipot.net/hazo/server/views"
)

//go:embed characters.html
var charactersPage string

type State struct {
	PageTitle         string
	DatasetName       string
	AvailableDatasets []db.Dataset
	MenuState         *treemenu.State
	ViewMenuState     *popover.State
}

func Handler(w http.ResponseWriter, r *http.Request, cc *common.Context) error {
	tmpl := components.NewTemplate()
	dbName := r.PathValue("dsName")
	if dbName == "" {
		dbName = "plants"
	}
	docId := r.PathValue("id")
	ctx := context.Background()
	queries, err := db.Open(fmt.Sprintf("%s.sq3", dbName))
	if err != nil {
		return err
	}
	tmpl.Funcs(template.FuncMap{
		"selectedDoc": func() string {
			return docId
		},
		"sortDocs": func(items []*treemenu.Item) []*treemenu.Item {
			slices.SortFunc(items, func(i, o *treemenu.Item) int {
				return i.Order - o.Order
			})
			return items
		},
		"documentUrl": func(taxon *treemenu.Item) string {
			return fmt.Sprintf("/ds/%s/characters/%s", dbName, taxon.Id)
		},
	})
	tmpl, err = tmpl.Parse(charactersPage)
	if err != nil {
		return err
	}
	items, err := treemenu.LoadItemFromDb(ctx, queries, "c0", [3]string{"FR", "EN", "CN"})
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
		MenuState: &treemenu.State{
			Selected: docId,
			Root:     items,
		},
		ViewMenuState: views.NewMenuState("Characters", dbName),
	})
	if err != nil {
		return err
	}
	return nil
}
