package views

import (
	"context"
	"fmt"
	"strings"

	"nicolas.galipot.net/hazo/db"
	"nicolas.galipot.net/hazo/db/storage"
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

type linkMaker = func(dsName string, ref string) string

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

func LinkToDescriptor(taxonRef string) linkMaker {
	return func(dsName string, ref string) string {
		return fmt.Sprintf("/ds/%s/taxons/%s?d=%s", dsName, taxonRef, ref)
	}
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

type Descriptor struct {
	Ref        string
	Name       string
	NameTr1    string
	NameTr2    string
	Url        string
	Source     string
	IsSelected bool
}

func GetTaxonDescriptors(ctx context.Context, queries *db.Queries, dsName string, taxonRef string, currentDescriptor *DocState) ([]Descriptor, error) {
	rows, err := queries.GetDescriptors(ctx, storage.GetDescriptorsParams{
		Path:     db.FullPath(currentDescriptor.Path, currentDescriptor.Ref),
		TaxonRef: taxonRef,
	})
	if err != nil {
		return nil, err
	}
	descriptors := make([]Descriptor, len(rows))
	for i, row := range rows {
		unsel, ok := row.Unselected.(int64)
		isSelected := false
		if ok {
			isSelected = unsel == 0
		}
		descriptors[i] = Descriptor{
			Ref:        row.Ref,
			Name:       row.Name,
			NameTr1:    row.NameTr1.String,
			NameTr2:    row.NameTr2.String,
			Url:        LinkToDescriptor(taxonRef)(dsName, row.Ref),
			Source:     row.Source.String,
			IsSelected: isSelected,
		}
	}
	return descriptors, nil
}

func GetDocumentBranch(ctx context.Context, queries *db.Queries, doc *DocState, dbName string, makeLink linkMaker) (*breadcrumbs.State, error) {
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
		model[i].Url = makeLink(dbName, doc.Ref)
	}
	model[len(model)-1] = breadcrumbs.BreadCrumb{
		Label: doc.Name,
		Url:   makeLink(dbName, doc.Ref),
	}
	return &breadcrumbs.State{Branch: model}, nil
}
