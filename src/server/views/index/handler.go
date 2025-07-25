package index

import (
	_ "embed"
	"fmt"
	"html/template"
	"net/http"

	"nicolas.galipot.net/hazo/server/common"
	"nicolas.galipot.net/hazo/server/components"
	"nicolas.galipot.net/hazo/storage"
)

//go:embed index.html
var indexPage string

type ViewModel struct {
	PageTitle              string
	Datasets               []storage.Dataset
	SharedReadableDatasets []storage.SharedDataset
	SharedWritableDatasets []storage.SharedDataset
	Debug                  bool
}

func Handler(w http.ResponseWriter, r *http.Request, cc *common.Context) error {
	tmpl := components.NewTemplate()
	tmpl = template.Must(tmpl.Parse(indexPage))
	datasets, err := cc.User.ListDatasets()
	if err != nil {
		return fmt.Errorf("failed to list datasets in index handler: %w", err)
	}
	sharedReadableDatasets, err := cc.User.GetReadableSharedDatasets()
	if err != nil {
		return fmt.Errorf("failed to list datasets in index handler: %w", err)
	}
	sharedWritableDatasets, err := cc.User.GetWritableSharedDatasets()
	if err != nil {
		return fmt.Errorf("failed to list datasets in index handler: %w", err)
	}
	w.Header().Add("Content-Type", "text/html")
	err = tmpl.Execute(w, ViewModel{
		PageTitle:              "Hazo Home",
		Datasets:               datasets,
		SharedReadableDatasets: sharedReadableDatasets,
		SharedWritableDatasets: sharedWritableDatasets,
		Debug:                  cc.Config.Debug,
	})
	if err != nil {
		return fmt.Errorf("template rendering of the index page failed: %w", err)
	}
	return nil
}
