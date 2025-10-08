package characters

import (
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"net/http"

	"nicolas.galipot.net/hazo/server/common"
	"nicolas.galipot.net/hazo/server/components/picturebox"
	"nicolas.galipot.net/hazo/server/components/popover"
	"nicolas.galipot.net/hazo/server/components/treemenu"
	"nicolas.galipot.net/hazo/server/documents"
	"nicolas.galipot.net/hazo/server/link"
	"nicolas.galipot.net/hazo/server/views"
	"nicolas.galipot.net/hazo/storage/dataset"
	"nicolas.galipot.net/hazo/storage/dsdb"
)

//go:embed characters.html
var charactersPage string

func Handler(w http.ResponseWriter, r *http.Request, cc *common.Context) error {
	dsName := r.PathValue("dsName")
	docRef := r.PathValue("id")
	ctx := context.Background()
	cc.ExecuteActions(ctx, r)
	ds, err := cc.User.GetDataset(dsName)
	if err != nil {
		return err
	}
	queries, err := dataset.OpenDb(ds)
	if err != nil {
		return err
	}
	queryParams := r.URL.Query()
	menuLangs, menuSelectedLangNames, err := documents.LoadMenuLanguages(ctx, cc)
	if err != nil {
		return fmt.Errorf("loading taxon languages: %w", err)
	}
	template.Must(cc.Template.Parse(charactersPage))
	items, err := treemenu.LoadItemFromDb(ctx, queries, "c0", menuSelectedLangNames, queryParams.Get("filterMenu"))
	if err != nil {
		return err
	}
	datasets, err := views.NewDatasetMenuViewModel(cc, dsName)
	if err != nil {
		return err
	}
	var character *documents.ViewModel
	ch, err := queries.GetDocumentTr2(ctx, dsdb.GetDocumentTr2Params{
		Lang1: "EN",
		Lang2: "CN",
		Ref:   docRef,
	})
	if err == nil {
		// TODO: handle non empty row error
		character = &documents.ViewModel{
			Ref:         ch.Ref,
			Path:        ch.Path,
			Name:        ch.Name,
			NameEN:      ch.NameTr1.String,
			NameCN:      ch.NameTr2.String,
			Description: ch.Details.String,
		}
	} else {
		fmt.Printf("error: %s\n", err.Error())
	}
	breadCrumbs, err := views.GetDocumentBranch(ctx, queries, character, dsName, link.ToCharacter)
	if err != nil {
		return err
	}
	attach, err := queries.GetDocumentAttachments(ctx, ch.Ref)
	picboxModel := picturebox.ViewModel{Index: 0, Count: 0, Name: ch.Name}
	if err == nil {
		picboxModel.Count = len(attach)
		if len(attach) > 0 {
			picboxModel.Index = 1
			picboxModel.Source = attach[0].Source
		}
	}
	template.Must(cc.Template.Parse(FormTemplate))
	err = cc.Template.Execute(w, ViewModel{
		PageTitle:         "Hazo",
		DatasetName:       dsName,
		Debug:             cc.Config.Debug,
		AvailableDatasets: datasets,
		LangsCheckList: popover.CheckListViewModel{
			Label: "",
			Icon:  "fa-language",
			Items: menuLangs,
		},
		MenuState: &treemenu.ViewModel{
			Selected:     docRef,
			ColumnsCount: len(menuSelectedLangNames),
			Root:         items,
		},
		MenuViewModel:     views.NewViewMenuViewModel("Characters", dsName),
		BreadCrumbsState:  breadCrumbs,
		SelectedCharacter: character,
		PictureBoxModel:   &picboxModel,
	})
	if err != nil {
		return err
	}
	return nil
}
