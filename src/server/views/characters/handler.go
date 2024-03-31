package characters

import (
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"net/http"

	"nicolas.galipot.net/hazo/db"
	"nicolas.galipot.net/hazo/db/storage"
	"nicolas.galipot.net/hazo/server/common"
	"nicolas.galipot.net/hazo/server/components/breadcrumbs"
	"nicolas.galipot.net/hazo/server/components/picturebox"
	"nicolas.galipot.net/hazo/server/components/popover"
	"nicolas.galipot.net/hazo/server/components/treemenu"
	"nicolas.galipot.net/hazo/server/views"
)

//go:embed characters.html
var charactersPage string

type State struct {
	PageTitle         string
	DatasetName       string
	AvailableDatasets *popover.State
	MenuState         *treemenu.State
	ViewMenuState     *popover.State
	BreadCrumbsState  *breadcrumbs.State
	SelectedCharacter *views.DocState
	PictureBoxModel   *picturebox.Model
}

func Handler(w http.ResponseWriter, r *http.Request, cc *common.Context) error {
	dbName := r.PathValue("dsName")
	docRef := r.PathValue("id")
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
	datasets, err := views.NewDatasetMenuState(dbName)
	if err != nil {
		return err
	}
	var character *views.DocState
	ch, err := queries.GetDocumentTr2(ctx, storage.GetDocumentTr2Params{
		Lang1: "EN",
		Lang2: "CN",
		Ref:   docRef,
	})
	if err == nil {
		// TODO: handle non empty row error
		character = &views.DocState{
			Ref:         ch.Ref,
			Path:        ch.Path,
			Name:        ch.Name,
			NameEN:      ch.NameTr1.String,
			NameCN:      ch.NameTr2.String,
			Description: ch.Details.String,
		}
	}
	breadCrumbs, err := views.GetDocumentBranch(ctx, queries, character, dbName, "characters")
	if err != nil {
		return err
	}
	attach, err := queries.GetDocumentAttachments(ctx, ch.Ref)
	picboxModel := picturebox.Model{Index: 0, Count: 0, Name: ch.Name}
	if err == nil {
		picboxModel.Count = len(attach)
		if len(attach) > 0 {
			picboxModel.Index = 1
			picboxModel.Source = attach[0].Source
		}
	}
	template.Must(cc.Template.Parse(FormTemplate))
	err = cc.Template.Execute(w, State{
		PageTitle:         "Hazo",
		DatasetName:       dbName,
		AvailableDatasets: datasets,
		MenuState: &treemenu.State{
			Selected: docRef,
			Root:     items,
		},
		ViewMenuState:     views.NewMenuState("Characters", dbName),
		BreadCrumbsState:  breadCrumbs,
		SelectedCharacter: character,
		PictureBoxModel:   &picboxModel,
	})
	if err != nil {
		return err
	}
	return nil
}
