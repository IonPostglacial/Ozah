package characters

import (
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"net/http"

	"nicolas.galipot.net/hazo/db"
	"nicolas.galipot.net/hazo/server/common"
	"nicolas.galipot.net/hazo/server/components/breadcrumbs"
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
	BreadCrumbsState  *breadcrumbs.State
}

func Handler(w http.ResponseWriter, r *http.Request, cc *common.Context) error {
	dbName := r.PathValue("dsName")
	docId := r.PathValue("id")
	ctx := context.Background()
	queries, err := db.Open(fmt.Sprintf("%s.sq3", dbName))
	if err != nil {
		return err
	}
	template.Must(cc.Template.Parse(charactersPage))
	items, err := treemenu.LoadItemFromDb(ctx, queries, "c0", [3]string{"FR", "EN", "CN"})
	if err != nil {
		return err
	}
	datasets, err := db.ListDatasets()
	if err != nil {
		return err
	}
	ch, err := queries.GetDocument(ctx, docId)
	if err != nil {
		return err
	}
	breadCrumbs, err := views.GetDocumentBranch(ctx, queries, ch.Path, dbName, "characters")
	if err != nil {
		return err
	}
	breadCrumbs.Branch = append(breadCrumbs.Branch, breadcrumbs.BreadCrumb{
		Label: ch.Name,
		Url:   fmt.Sprintf("/ds/%s/characters/%s", dbName, docId),
	})
	err = cc.Template.Execute(w, State{
		PageTitle:         "Hazo",
		DatasetName:       dbName,
		AvailableDatasets: datasets,
		MenuState: &treemenu.State{
			Selected: docId,
			Root:     items,
		},
		ViewMenuState:    views.NewMenuState("Characters", dbName),
		BreadCrumbsState: breadCrumbs,
	})
	if err != nil {
		return err
	}
	return nil
}
