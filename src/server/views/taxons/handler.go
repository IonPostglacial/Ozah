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

type Model struct {
	PageTitle                   string
	DatasetName                 string
	AvailableDatasets           *popover.State
	MenuState                   *treemenu.State
	SelectedTaxon               *FormData
	ViewMenuState               *popover.State
	BreadCrumbsState            *breadcrumbs.State
	DescriptorsBreadCrumbsState *breadcrumbs.State
	Descriptors                 []views.Descriptor
	SummaryModel                *summary.Model
	PictureBoxModel             *picturebox.Model
	BookInfoModel               []storage.GetTaxonBookInfoRow
}

func Handler(w http.ResponseWriter, r *http.Request, cc *common.Context) error {
	dsName := r.PathValue("dsName")
	docRef := r.PathValue("id")
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
	if err != nil {
		return err
	}
	queryParams := r.URL.Query()
	descriptorRef := queryParams.Get("d")
	var currentDescriptor *views.DocState
	if descriptorRef == "" {
		descriptorRef = "c0"
		currentDescriptor = &views.DocState{Ref: descriptorRef, Path: ""}
	} else {
		doc, err := queries.GetDocument(ctx, descriptorRef)
		if err != nil {
			return err
		}
		currentDescriptor = &views.DocState{Ref: doc.Ref, Path: doc.Path, Name: doc.Name}
	}
	if docRef != "" {
		taxon, err = LoadFormDataFromDb(ctx, queries, docRef)
		if err != nil {
			return err
		}
	} else {
		taxon = &FormData{}
	}
	template.Must(cc.Template.Parse(taxonPage))
	template.Must(cc.Template.Parse(FormTemplate))
	items, err := treemenu.LoadItemFromDb(ctx, queries, "t0", [3]string{"S", "V", "CN"}, queryParams.Get("filterMenu"))
	if err != nil {
		return err
	}
	datasets, err := views.NewDatasetMenuState(cc, dsName)
	if err != nil {
		return err
	}
	branch, err := views.GetDocumentBranch(ctx, queries, &taxon.DocState, dsName, views.LinkToTaxon)
	if err != nil {
		return err
	}
	descBreadcrumbs, err := views.GetDocumentBranch(ctx, queries, currentDescriptor, dsName, views.LinkToDescriptor(docRef))
	if err != nil {
		return err
	}
	// TODO: retrieve selection by taxon
	descriptors, err := views.GetTaxonDescriptors(ctx, queries, dsName, taxon.Ref, currentDescriptor)
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
	bookInfo, err := queries.GetTaxonBookInfo(ctx, taxon.Ref)
	if err != nil {
		return err
	}
	err = cc.Template.Execute(w, Model{
		PageTitle:         "Hazo",
		DatasetName:       dsName,
		AvailableDatasets: datasets,
		SelectedTaxon:     taxon,
		MenuState: &treemenu.State{
			Selected: taxon.Ref,
			Langs:    []string{"S", "V", "CN"},
			Root:     items,
		},
		ViewMenuState:               views.NewMenuState("Taxons", dsName),
		BreadCrumbsState:            branch,
		DescriptorsBreadCrumbsState: descBreadcrumbs,
		Descriptors:                 descriptors,
		SummaryModel:                summary,
		PictureBoxModel:             &picboxModel,
		BookInfoModel:               bookInfo,
	})
	if err != nil {
		return err
	}
	return nil
}
