package export

import (
	"fmt"
	"net/http"

	"nicolas.galipot.net/hazo/server/common"
	"nicolas.galipot.net/hazo/storage/dataset"
)

func JsonHandler(w http.ResponseWriter, r *http.Request, cc *common.Context) error {
	dsName := r.PathValue("dsName")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	pds, err := cc.User.GetDataset(dsName)
	if err != nil {
		return fmt.Errorf("could not get dataset '%s': %w", dsName, err)
	}
	queries, err := dataset.OpenDb(pds)
	if err != nil {
		return fmt.Errorf("could not open dataset database for '%s': %w", dsName, err)
	}
	err = dataset.ExportJson(dsName, queries, w)
	if err != nil {
		return fmt.Errorf("failed to export dataset as Hazo JSON: %w", err)
	}
	return nil
}
