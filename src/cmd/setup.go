package cmd

import (
	"context"
	"fmt"

	"nicolas.galipot.net/hazo/storage/app"
	"nicolas.galipot.net/hazo/storage/appdb"
)

func Setup(args []string) error {
	db, queries, err := app.OpenDb()
	if err != nil {
		return fmt.Errorf("couldn't open appdb: %w", err)
	}
	ctx := context.Background()
	_, err = db.Exec(app.Schema)
	if err != nil {
		return fmt.Errorf("couldn't apply database schema during setup: %w", err)
	}
	langs := []appdb.InsertLangParams{
		{Ref: "V", Name: "Vernacular"},
		{Ref: "CN", Name: "Chinese"},
		{Ref: "EN", Name: "English"},
		{Ref: "FR", Name: "French"},
	}
	for _, lang := range langs {
		_, err = queries.InsertLang(ctx, lang)
		if err != nil {
			return fmt.Errorf("could not insert lang during setup: %w", err)
		}
	}
	panels := []appdb.InsertUserPanelParams{
		{ID: 0, Name: "Properties"},
		{ID: 1, Name: "Descriptors"},
		{ID: 2, Name: "Summary"},
	}
	for _, panel := range panels {
		_, err = queries.InsertUserPanel(ctx, panel)
		if err != nil {
			return fmt.Errorf("could not insert panel during setup: %w", err)
		}
	}
	if err != nil {
		return fmt.Errorf("couldn't apply appdb schema: %w", err)
	}
	return nil
}
