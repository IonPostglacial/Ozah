package iconmenu

import (
	"context"

	"nicolas.galipot.net/hazo/db"
	"nicolas.galipot.net/hazo/db/storage"
	"nicolas.galipot.net/hazo/server/documents"
	"nicolas.galipot.net/hazo/server/link"
)

type Model struct {
	Ref        string
	Name       string
	NameTr1    string
	NameTr2    string
	Url        string
	Source     string
	IsSelected bool
}

func GetTaxonDescriptors(ctx context.Context, queries *db.Queries, dsName string, taxonRef string, currentDescriptor *documents.Model) ([]Model, error) {
	rows, err := queries.GetDescriptors(ctx, storage.GetDescriptorsParams{
		Path:     db.FullPath(currentDescriptor.Path, currentDescriptor.Ref),
		TaxonRef: taxonRef,
	})
	if err != nil {
		return nil, err
	}
	descriptors := make([]Model, len(rows))
	for i, row := range rows {
		unsel, ok := row.Unselected.(int64)
		isSelected := false
		if ok {
			isSelected = unsel == 0
		}
		descriptors[i] = Model{
			Ref:        row.Ref,
			Name:       row.Name,
			NameTr1:    row.NameTr1.String,
			NameTr2:    row.NameTr2.String,
			Url:        link.ToDescriptor(taxonRef)(dsName, row.Ref),
			Source:     row.Source.String,
			IsSelected: isSelected,
		}
	}
	return descriptors, nil
}
