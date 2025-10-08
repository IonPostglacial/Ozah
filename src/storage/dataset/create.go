package dataset

import "fmt"

func Create(ds Private) error {
	dtb, err := ConnectDb(ds)
	if err != nil {
		return fmt.Errorf("connecting to database '%s' failed: %w", ds, err)
	}
	_, err = dtb.Exec(Schema)
	if err != nil {
		return fmt.Errorf("creating database scheme from '%s' failed: %w", Schema, err)
	}
	_, err = dtb.Exec(Index)
	if err != nil {
		return fmt.Errorf("creating database indices failed: %w", err)
	}
	return nil
}
