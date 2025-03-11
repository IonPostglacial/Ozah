package cmd

import (
	"fmt"

	"nicolas.galipot.net/hazo/storage"
)

func Setup(args []string) error {
	db, _, err := storage.OpenAppDb()
	if err != nil {
		return fmt.Errorf("Couldn't open appdb: %w", err)
	}
	_, err = db.Exec(storage.AppSchema)
	if err != nil {
		return fmt.Errorf("Couldn't apply appdb schema: %w", err)
	}
	return nil
}
