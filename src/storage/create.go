package storage

import "fmt"

func Create(ds PrivateDataset) error {
	dtb, err := ConnectDsDb(ds)
	if err != nil {
		return fmt.Errorf("connecting to database '%s' failed: %w", ds, err)
	}
	_, err = dtb.Exec(DatasetSchema)
	if err != nil {
		return fmt.Errorf("creating database scheme failed: %w", err)
	}
	_, err = dtb.Exec(index)
	if err != nil {
		return fmt.Errorf("creating database indices failed: %w", err)
	}
	return nil
}
