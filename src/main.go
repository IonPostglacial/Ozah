package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"nicolas.galipot.net/hazo/storage"
)

func connectToDatabase() (*storage.Queries, error) {
	db, err := sql.Open("sqlite3", "file:hello.sqlite")
	if err != nil {
		return nil, fmt.Errorf("could not connect to the database: %w", err)
	}
	queries := storage.New(db)
	return queries, nil
}

func insertStandardContent(ctx context.Context) error {
	queries, err := connectToDatabase()
	if err != nil {
		return fmt.Errorf("connection to databse to insert standard content failed: %w", err)
	}
	_, err = queries.InsertStdLangs(ctx)
	if err != nil {
		return fmt.Errorf("could not insert standard language types: %w", err)
	}
	return nil
}

func run(ctx context.Context) error {
	queries, err := connectToDatabase()
	if err != nil {
		return err
	}
	acanthaceae, err := queries.GetDocumentHierarchy(ctx, "t0.t1")
	if err != nil {
		return err
	}
	fmt.Printf("%#v\n", acanthaceae)
	return nil
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Pass command 'init' or 'run'.")
	}
	var err error
	ctx := context.Background()
	switch os.Args[1] {
	case "init":
		err = insertStandardContent(ctx)
	case "run":
		err = run(ctx)
	default:
		log.Fatal("unknown command:", os.Args[1])
	}
	if err != nil {
		log.Fatal(err)
	}
}
