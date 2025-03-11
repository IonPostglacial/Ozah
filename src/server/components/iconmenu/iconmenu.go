package iconmenu

import (
	"context"

	"nicolas.galipot.net/hazo/server/documents"
	"nicolas.galipot.net/hazo/server/link"
	"nicolas.galipot.net/hazo/storage"
	"nicolas.galipot.net/hazo/storage/dsdb"
)

type ViewModel struct {
	Ref        string
	Name       string
	NameTr1    string
	NameTr2    string
	Url        string
	Source     string
	IsSelected bool
}

func GetTaxonDescriptors(ctx context.Context, queries *storage.Queries, dsName string, taxonRef string, currentDescriptor *documents.ViewModel) ([]ViewModel, error) {
	rows, err := queries.GetDescriptors(ctx, dsdb.GetDescriptorsParams{
		Path:     storage.FullPath(currentDescriptor.Path, currentDescriptor.Ref),
		TaxonRef: taxonRef,
	})
	if err != nil {
		return nil, err
	}
	descriptors := make([]ViewModel, len(rows))
	for i, row := range rows {
		unsel, ok := row.Unselected.(int64)
		isSelected := false
		if ok {
			isSelected = unsel == 0
		}
		descriptors[i] = ViewModel{
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
