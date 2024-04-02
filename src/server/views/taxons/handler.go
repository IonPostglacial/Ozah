package taxons

import (
	"context"
	_ "embed"
	"html/template"
	"net/http"

	"nicolas.galipot.net/hazo/db"
	"nicolas.galipot.net/hazo/db/storage"
	"nicolas.galipot.net/hazo/server/common"
	"nicolas.galipot.net/hazo/server/components/breadcrumbs"
	"nicolas.galipot.net/hazo/server/components/picturebox"
	"nicolas.galipot.net/hazo/server/components/popover"
	"nicolas.galipot.net/hazo/server/components/summary"
	"nicolas.galipot.net/hazo/server/components/treemenu"
	"nicolas.galipot.net/hazo/server/views"
)

//go:embed taxons.html
var taxonPage string

type State struct {
	PageTitle                   string
	DatasetName                 string
	AvailableDatasets           *popover.State
	MenuState                   *treemenu.State
	SelectedTaxon               *FormData
	ViewMenuState               *popover.State
	BreadCrumbsState            *breadcrumbs.State
	DescriptorsBreadCrumbsState *breadcrumbs.State
	Descriptors                 []storage.GetDescriptorsRow
	SummaryModel                *summary.Model
	PictureBoxModel             *picturebox.Model
}

func Handler(w http.ResponseWriter, r *http.Request, cc *common.Context) error {
	dsName := r.PathValue("dsName")
	docId := r.PathValue("id")
	var (
		taxon *FormData
		err   error
	)
	ctx := context.Background()
	ds, err := cc.User.GetDataset(dsName)
	if err != nil {
		return err
	}
	queries, err := db.Open(ds)
	currentDescriptor := &views.DocState{Ref: "c0", Path: "c0"}
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
	datasets, err := views.NewDatasetMenuState(cc, dsName)
	if err != nil {
		return err
	}
	branch, err := views.GetDocumentBranch(ctx, queries, &taxon.DocState, dsName)
	if err != nil {
		return err
	}
	descBreadcrumbs, err := views.GetDocumentBranch(ctx, queries, currentDescriptor, dsName)
	if err != nil {
		return err
	}
	// TODO: retrieve selection by taxon
	descriptors, err := queries.GetDescriptors(ctx, storage.GetDescriptorsParams{
		Path:     currentDescriptor.Ref,
		TaxonRef: taxon.Ref,
	})
	if err != nil {
		return err
	}
	summary, err := summary.LoadForTaxon(ctx, queries, taxon.Ref)
	if err != nil {
		return err
	}
	attach, err := queries.GetDocumentAttachments(ctx, taxon.Ref)
	picboxModel := picturebox.Model{Index: 0, Count: 0, Name: taxon.Name}
	if err == nil {
		picboxModel.Count = len(attach)
		if len(attach) > 0 {
			picboxModel.Index = 1
			picboxModel.Source = attach[0].Source
		}
	}
	err = cc.Template.Execute(w, State{
		PageTitle:         "Hazo",
		DatasetName:       dsName,
		AvailableDatasets: datasets,
		SelectedTaxon:     taxon,
		MenuState: &treemenu.State{
			Selected: taxon.Ref,
			Root:     items,
		},
		ViewMenuState:               views.NewMenuState("Taxons", dsName),
		BreadCrumbsState:            branch,
		DescriptorsBreadCrumbsState: descBreadcrumbs,
		Descriptors:                 descriptors,
		SummaryModel:                summary,
		PictureBoxModel:             &picboxModel,
	})
	if err != nil {
		return err
	}
	return nil
}
