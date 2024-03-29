package views

import (
	"context"
	"fmt"
	"strings"

	"nicolas.galipot.net/hazo/db"
	"nicolas.galipot.net/hazo/db/storage"
	"nicolas.galipot.net/hazo/server/components/breadcrumbs"
	"nicolas.galipot.net/hazo/server/components/popover"
)

type DocState struct {
	Ref         string
	Path        string
	Name        string
	Description string
}

func NewMenuState(label, dsName string) *popover.State {
	return &popover.State{
		Label: label,
		Items: []popover.Item{
			{Url: fmt.Sprintf("/ds/%s/taxons", dsName), Label: "Taxons"},
			{Url: fmt.Sprintf("/ds/%s/characters", dsName), Label: "Characters"},
		},
	}
}

func NewDatasetMenuState(label string) (*popover.State, error) {
	datasets, err := db.ListDatasets()
	if err != nil {
		return nil, err
	}
	items := make([]popover.Item, len(datasets))
	for i, ds := range datasets {
		items[i] = popover.Item{
			Url:   fmt.Sprintf("/ds/%s/taxons", ds.Name),
			Label: ds.Name,
		}
	}
	return &popover.State{
		Label: label,
		Items: items,
	}, nil
}

func GetDocumentBranch(ctx context.Context, queries *storage.Queries, doc *DocState, dbName string, docType string) (*breadcrumbs.State, error) {
	if doc.Path == "" {
		return &breadcrumbs.State{}, nil
	}
	branch := strings.Split(doc.Path, ".")
	docs, err := queries.GetDocumentsNames(ctx, branch)
	if err != nil {
		return nil, err
	}
	model := make([]breadcrumbs.BreadCrumb, len(docs)+1)
	for i, doc := range docs {
		model[i].Label = doc.Name
		model[i].Url = fmt.Sprintf("/ds/%s/%s/%s", dbName, docType, doc.Ref)
	}
	model[len(model)-1] = breadcrumbs.BreadCrumb{
		Label: doc.Name,
		Url:   fmt.Sprintf("/ds/%s/%s/%s", dbName, docType, doc.Ref),
	}
	return &breadcrumbs.State{Branch: model}, nil
}
