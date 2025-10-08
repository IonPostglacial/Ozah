package taxons

import (
	"context"

	"nicolas.galipot.net/hazo/server/documents"
	"nicolas.galipot.net/hazo/storage/dataset"

	_ "embed"
)

//go:embed form.html
var FormTemplate string

func LoadFormViewModelFromDb(ctx context.Context, queries *dataset.Queries, id string) (*FormViewModel, error) {
	data, err := queries.GetTaxonInfo(ctx, id)
	if err != nil {
		return nil, err
	}
	return &FormViewModel{
		ViewModel: documents.ViewModel{
			Ref:         id,
			Path:        data.Path,
			Name:        data.Name,
			Description: data.Details.String,
		},
		NameV:   data.NameV.String,
		NameCN:  data.NameCn.String,
		Author:  data.Author,
		Website: data.Website.String,
	}, nil
}
