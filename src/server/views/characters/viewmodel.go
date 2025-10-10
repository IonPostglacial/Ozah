package characters

import (
	"nicolas.galipot.net/hazo/server/components/breadcrumbs"
	"nicolas.galipot.net/hazo/server/components/picturebox"
	"nicolas.galipot.net/hazo/server/components/popover"
	"nicolas.galipot.net/hazo/server/components/treemenu"
	"nicolas.galipot.net/hazo/server/documents"
)

type CharacterViewModel struct {
	documents.ViewModel
	Color string
}

type ViewModel struct {
	PageTitle         string
	DatasetName       string
	Debug             bool
	AvailableDatasets *popover.ViewModel
	LangsCheckList    popover.CheckListViewModel
	MenuState         *treemenu.ViewModel
	MenuViewModel     *popover.ViewModel
	BreadCrumbsState  *breadcrumbs.ViewModel
	SelectedCharacter *CharacterViewModel
	PictureBoxModel   *picturebox.ViewModel
}
