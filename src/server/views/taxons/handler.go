package taxons

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

//go:embed taxons.html
var taxonPage string

type State struct {
	PageTitle         string
	DatasetName       string
	AvailableDatasets []db.Dataset
	MenuState         *treemenu.State
	SelectedTaxon     *FormData
	ViewMenuState     *popover.State
	BreadCrumbsState  *breadcrumbs.State
}

func Handler(w http.ResponseWriter, r *http.Request, cc *common.Context) error {
	dbName := r.PathValue("dsName")
	docId := r.PathValue("id")
	var (
		taxon *FormData
		err   error
	)
	ctx := context.Background()
	queries, err := db.Open(fmt.Sprintf("%s.sq3", dbName))
	if err != nil {
		return err
	}
	if docId != "" {
		taxon, err = LoadFormDataFromDb(ctx, queries, docId)
		if err != nil {
			return err
		}
	} else {
		taxon = &FormData{}
	}
	template.Must(cc.Template.Parse(taxonPage))
	template.Must(cc.Template.Parse(FormTemplate))
	items, err := treemenu.LoadItemFromDb(ctx, queries, "t0", [3]string{"S", "V", "CN"})
	if err != nil {
		return err
	}
	datasets, err := db.ListDatasets()
	if err != nil {
		return err
	}
	branch, err := views.GetDocumentBranch(ctx, queries, &taxon.DocState, dbName, "taxons")
	if err != nil {
		return err
	}
	err = cc.Template.Execute(w, State{
		PageTitle:         "Hazo",
		DatasetName:       dbName,
		AvailableDatasets: datasets,
		SelectedTaxon:     taxon,
		MenuState: &treemenu.State{
			Selected: taxon.Id,
			Root:     items,
		},
		ViewMenuState:    views.NewMenuState("Taxons", dbName),
		BreadCrumbsState: branch,
	})
	if err != nil {
		return err
	}
	return nil
}
