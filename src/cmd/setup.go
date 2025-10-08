package cmd

import (
	"context"
	"flag"
	"fmt"

	"nicolas.galipot.net/hazo/storage/app"
	"nicolas.galipot.net/hazo/storage/appdb"
)

func Setup(args []string) error {
	fs := flag.NewFlagSet("setup", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: hazo setup\n\n")
		fmt.Fprintf(fs.Output(), "Initialize the application database with default configuration.\n\n")
		fmt.Fprintf(fs.Output(), "This command creates the necessary tables and populates them with:\n")
		fmt.Fprintf(fs.Output(), "  - Default languages (Vernacular, Chinese, English, French)\n")
		fmt.Fprintf(fs.Output(), "  - Default user panels (Properties, Descriptors, Summary)\n\n")
		fmt.Fprintf(fs.Output(), "Options:\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return err
	}
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
