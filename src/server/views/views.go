package views

import (
	"context"
	"fmt"
	"strings"

	"nicolas.galipot.net/hazo/db"
	"nicolas.galipot.net/hazo/server/common"
	"nicolas.galipot.net/hazo/server/components/breadcrumbs"
	"nicolas.galipot.net/hazo/server/components/popover"
)

type DocState struct {
	Ref         string
	Path        string
	Name        string
	NameEN      string
	NameCN      string
	Description string
}

func LinkToTaxons(dsName string) string {
	return fmt.Sprintf("/ds/%s/taxons", dsName)
}

func LinkToTaxon(dsName string, ref string) string {
	return fmt.Sprintf("/ds/%s/taxons/%s", dsName, ref)
}

func LinkToCharacters(dsName string) string {
	return fmt.Sprintf("/ds/%s/characters", dsName)
}

func LinkToCharacter(dsName string, ref string) string {
	return fmt.Sprintf("/ds/%s/characters/%s", dsName, ref)
}

func LinkToIdentify(dsName string) string {
	return fmt.Sprintf("/ds/%s/identify", dsName)
}

func LinkToDocument(dsName string, ref string) string {
	switch {
	case strings.HasPrefix(ref, "t"):
		return LinkToTaxon(dsName, ref)
	case strings.HasPrefix(ref, "c"):
		return LinkToCharacter(dsName, ref)
	default:
		return LinkToTaxons(dsName)
	}
}

func NewMenuState(label, dsName string) *popover.State {
	return &popover.State{
		Label: label,
		Items: []popover.Item{
			{Url: LinkToTaxons(dsName), Label: "Taxons"},
			{Url: LinkToCharacters(dsName), Label: "Characters"},
			{Url: LinkToIdentify(dsName), Label: "Identification"},
		},
	}
}

func NewDatasetMenuState(cc *common.Context, label string) (*popover.State, error) {
	datasets, err := cc.User.ListDatasets()
	if err != nil {
		return nil, err
	}
	items := make([]popover.Item, len(datasets))
	for i, ds := range datasets {
		items[i] = popover.Item{
			Url:   LinkToTaxons(ds.Name),
			Label: ds.Name,
		}
	}
	return &popover.State{
		Label: label,
		Items: items,
	}, nil
}

func GetDocumentBranch(ctx context.Context, queries *db.Queries, doc *DocState, dbName string) (*breadcrumbs.State, error) {
	if doc == nil || doc.Path == "" {
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
		model[i].Url = LinkToDocument(dbName, doc.Ref)
	}
	model[len(model)-1] = breadcrumbs.BreadCrumb{
		Label: doc.Name,
		Url:   LinkToDocument(dbName, doc.Ref),
	}
	return &breadcrumbs.State{Branch: model}, nil
}
