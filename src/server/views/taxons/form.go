package taxons

import (
	"context"

	"nicolas.galipot.net/hazo/db/storage"

	_ "embed"
)

//go:embed form.html
var FormTemplate string

type FormData struct {
	Id          string
	Name        string
	NameV       string
	NameCN      string
	Author      string
	Website     string
	Description string
}

func LoadFormDataFromDb(ctx context.Context, queries *storage.Queries, id string) (*FormData, error) {
	data, err := queries.GetTaxonInfo(ctx, id)
	if err != nil {
		return nil, err
	}
	return &FormData{
		Id:          id,
		Name:        data.Name,
		NameV:       data.NameV.String,
		NameCN:      data.NameCn.String,
		Description: data.Details.String,
		Author:      data.Author,
		Website:     data.Website.String,
	}, nil
}
