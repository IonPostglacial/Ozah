package views

import (
	"context"
	"fmt"
	"strings"

	"nicolas.galipot.net/hazo/db/storage"
	"nicolas.galipot.net/hazo/server/components/breadcrumbs"
	"nicolas.galipot.net/hazo/server/components/popover"
)

func NewMenuState(label, dsName string) *popover.State {
	return &popover.State{
		Label: label,
		Items: []popover.Item{
			{Url: fmt.Sprintf("/ds/%s/taxons", dsName), Label: "Taxons"},
			{Url: fmt.Sprintf("/ds/%s/characters", dsName), Label: "Characters"},
		},
	}
}

func GetDocumentBranch(ctx context.Context, queries *storage.Queries, path string, dbName string, docType string) (*breadcrumbs.State, error) {
	branch := strings.Split(path, ".")
	names, err := queries.GetDocumentsNames(ctx, branch)
	if err != nil {
		return nil, err
	}
	model := make([]breadcrumbs.BreadCrumb, len(names))
	if len(names) < 1 {
		return &breadcrumbs.State{Branch: model}, nil
	}
	for i, name := range names {
		model[i].Label = name
		model[i].Url = fmt.Sprintf("/ds/%s/%s/%s", dbName, docType, branch[i])
	}
	return &breadcrumbs.State{Branch: model}, nil
}
