package taxons

import (
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"net/http"

	"slices"

	"nicolas.galipot.net/hazo/server/common"
	"nicolas.galipot.net/hazo/server/components/iconmenu"
	"nicolas.galipot.net/hazo/server/components/picturebox"
	"nicolas.galipot.net/hazo/server/components/summary"
	"nicolas.galipot.net/hazo/server/components/treemenu"
	"nicolas.galipot.net/hazo/server/documents"
	"nicolas.galipot.net/hazo/server/link"
	"nicolas.galipot.net/hazo/server/views"
	"nicolas.galipot.net/hazo/storage"
)

//go:embed taxons.html
var taxonPage string

func Handler(w http.ResponseWriter, r *http.Request, cc *common.Context) error {
	dsName := r.PathValue("dsName")
	docRef := r.PathValue("id")
	var (
		taxon *FormViewModel
		err   error
	)
	ctx := context.Background()
	ds, err := cc.User.GetDataset(dsName)
	if err != nil {
		return err
	}
	queries, err := storage.OpenDsDb(ds)
	if err != nil {
		return err
	}
	queryParams := r.URL.Query()
	menuLangSet := treemenu.LangSetFromString(queryParams.Get("menuLangs"))
	menuLangNames := []string{"S", "V", "CN"}
	menuSelectedLangs := menuLangSet.MaskNames(menuLangNames)
	menuLangs := menuLangSet.LangsFromNames(r.URL, menuLangNames)
	selectedPanels := PanelSetFromString(queryParams.Get("panels"))
	selectedPanelNames, unselectedPanels := selectedPanels.DivideNamesByMask(panelNames)
	descriptorRef := queryParams.Get("d")
	var currentDescriptor *documents.ViewModel
	if descriptorRef == "" {
		descriptorRef = "c0"
		currentDescriptor = &documents.ViewModel{Ref: descriptorRef, Path: ""}
	} else {
		doc, err := queries.GetDocument(ctx, descriptorRef)
		if err != nil {
			return err
		}
		currentDescriptor = &documents.ViewModel{Ref: doc.Ref, Path: doc.Path, Name: doc.Name}
	}
	if docRef != "" {
		taxon, err = LoadFormViewModelFromDb(ctx, queries, docRef)
		if err != nil {
			return err
		}
	} else {
		taxon = &FormViewModel{}
	}
	cc.Template.Funcs(template.FuncMap{
		"isPanelVisible": func(panelName string) bool {
			return slices.Contains(selectedPanelNames, panelName)
		},
		"panelZoomUrl": func(panel Panel) string {
			return PanelSet{common.BitSet(panel)}.LinkToPanelState(r.URL)
		},
		"panelAddUrl": func(panel uint64) string {
			return PanelSet{selectedPanels.With(common.BitSet(panel))}.LinkToPanelState(r.URL)
		},
		"panelRemoveUrl": func(panel Panel) string {
			return PanelSet{selectedPanels.Without(common.BitSet(panel))}.LinkToPanelState(r.URL)
		},
	})
	template.Must(cc.Template.Parse(taxonPage))
	template.Must(cc.Template.Parse(FormTemplate))
	items, err := treemenu.LoadItemFromDb(ctx, queries, "t0", menuSelectedLangs, queryParams.Get("filterMenu"))
	if err != nil {
		return err
	}
	datasets, err := views.NewDatasetMenuViewModel(cc, dsName)
	if err != nil {
		return err
	}
	branch, err := views.GetDocumentBranch(ctx, queries, &taxon.ViewModel, dsName, link.ToTaxon)
	if err != nil {
		return err
	}
	descBreadcrumbs, err := views.GetDocumentBranch(ctx, queries, currentDescriptor, dsName, link.ToDescriptor(docRef))
	if err != nil {
		return err
	}
	// TODO: retrieve selection by taxon
	descriptors, err := iconmenu.GetTaxonDescriptors(ctx, queries, dsName, taxon.Ref, currentDescriptor)
	if err != nil {
		return err
	}
	summary, err := summary.LoadForTaxon(ctx, queries, taxon.Ref)
	if err != nil {
		return err
	}
	attach, err := queries.GetDocumentAttachments(ctx, taxon.Ref)
	picboxModel := picturebox.ViewModel{Index: 0, Count: 0, Name: taxon.Name}
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
	err = cc.Template.Execute(w, ViewModel{
		PageTitle:         "Hazo",
		DatasetName:       dsName,
		Debug:             cc.Config.Debug,
		AvailableDatasets: datasets,
		SelectedTaxon:     taxon,
		MenuState: &treemenu.ViewModel{
			Selected:     taxon.Ref,
			Langs:        menuLangs,
			ColumnsCount: len(menuSelectedLangs),
			Root:         items,
		},
		MenuViewModel:               views.NewViewMenuViewModel("Taxons", dsName),
		BreadCrumbsState:            branch,
		DescriptorsBreadCrumbsState: descBreadcrumbs,
		Descriptors:                 descriptors,
		SummaryModel:                summary,
		PictureBoxModel:             &picboxModel,
		BookInfoModel:               bookInfo,
		UnselectedPanels:            unselectedPanels,
	})
	if err != nil {
		return fmt.Errorf("taxons template execution failed: %w", err)
	}
	return nil
}
