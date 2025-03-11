package views

import (
	"context"
	"strings"

	"nicolas.galipot.net/hazo/server/common"
	"nicolas.galipot.net/hazo/server/components/breadcrumbs"
	"nicolas.galipot.net/hazo/server/components/popover"
	"nicolas.galipot.net/hazo/server/documents"
	"nicolas.galipot.net/hazo/server/link"
	"nicolas.galipot.net/hazo/storage"
)

func NewViewMenuViewModel(label, dsName string) *popover.ViewModel {
	return &popover.ViewModel{
		Label: label,
		Items: []popover.Item{
			{Url: link.ToTaxons(dsName), Label: "Taxons"},
			{Url: link.ToCharacters(dsName), Label: "Characters"},
			{Url: link.ToIdentify(dsName), Label: "Identification"},
		},
	}
}

func NewDatasetMenuViewModel(cc *common.Context, label string) (*popover.ViewModel, error) {
	datasets, err := cc.User.ListDatasets()
	if err != nil {
		return nil, err
	}
	items := make([]popover.Item, len(datasets))
	for i, ds := range datasets {
		items[i] = popover.Item{
			Url:   link.ToTaxons(ds.Name),
			Label: ds.Name,
		}
	}
	return &popover.ViewModel{
		Label: label,
		Items: items,
	}, nil
}

func GetDocumentBranch(ctx context.Context, queries *storage.Queries, doc *documents.ViewModel, dbName string, makeLink link.Maker) (*breadcrumbs.ViewModel, error) {
	if doc == nil || doc.Path == "" {
		return &breadcrumbs.ViewModel{}, nil
	}
	branch := strings.Split(doc.Path, ".")
	docs, err := queries.GetDocumentsNames(ctx, branch)
	if err != nil {
		return nil, err
	}
	model := make([]breadcrumbs.BreadCrumb, len(docs)+1)
	for i, doc := range docs {
		model[i].Label = doc.Name
		model[i].Url = makeLink(dbName, doc.Ref)
	}
	model[len(model)-1] = breadcrumbs.BreadCrumb{
		Label: doc.Name,
		Url:   makeLink(dbName, doc.Ref),
	}
	return &breadcrumbs.ViewModel{Branch: model}, nil
}
