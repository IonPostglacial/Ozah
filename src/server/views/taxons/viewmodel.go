package taxons

import (
	"nicolas.galipot.net/hazo/db/storage"
	"nicolas.galipot.net/hazo/server/common"
	"nicolas.galipot.net/hazo/server/components/breadcrumbs"
	"nicolas.galipot.net/hazo/server/components/iconmenu"
	"nicolas.galipot.net/hazo/server/components/picturebox"
	"nicolas.galipot.net/hazo/server/components/popover"
	"nicolas.galipot.net/hazo/server/components/summary"
	"nicolas.galipot.net/hazo/server/components/treemenu"
	"nicolas.galipot.net/hazo/server/documents"
)

type FormViewModel struct {
	documents.ViewModel
	NameV   string
	NameCN  string
	Author  string
	Website string
}

type ViewModel struct {
	PageTitle                   string
	DatasetName                 string
	Debug                       bool
	AvailableDatasets           *popover.ViewModel
	MenuState                   *treemenu.ViewModel
	SelectedTaxon               *FormViewModel
	MenuViewModel               *popover.ViewModel
	BreadCrumbsState            *breadcrumbs.ViewModel
	DescriptorsBreadCrumbsState *breadcrumbs.ViewModel
	Descriptors                 []iconmenu.ViewModel
	SummaryModel                *summary.ViewModel
	PictureBoxModel             *picturebox.ViewModel
	BookInfoModel               []storage.GetTaxonBookInfoRow
	UnselectedPanels            []common.UnselectedItem
}
