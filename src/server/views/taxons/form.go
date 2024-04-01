package taxons

import (
	"context"

	"nicolas.galipot.net/hazo/db"
	"nicolas.galipot.net/hazo/server/views"

	_ "embed"
)

//go:embed form.html
var FormTemplate string

type FormData struct {
	views.DocState
	NameV   string
	NameCN  string
	Author  string
	Website string
}

func LoadFormDataFromDb(ctx context.Context, queries *db.Queries, id string) (*FormData, error) {
	data, err := queries.GetTaxonInfo(ctx, id)
	if err != nil {
		return nil, err
	}
	return &FormData{
		DocState: views.DocState{
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
