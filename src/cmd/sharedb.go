package cmd

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"nicolas.galipot.net/hazo/storage/app"
	"nicolas.galipot.net/hazo/storage/appdb"
)

var (
	ErrInvalidArgs = errors.New("invalid arguments")
	ErrInvalidMode = errors.New("invalid mode, must be 'read' or 'write'")
)

func Sharedb(args []string) error {
	fs := flag.NewFlagSet("sharedb", flag.ExitOnError)

	var creator, mode, datasetName, users string
	fs.StringVar(&creator, "creator", "", "Username of the dataset creator (required)")
	fs.StringVar(&mode, "mode", "", "Access mode: 'read' or 'write' (required)")
	fs.StringVar(&datasetName, "dataset", "", "Name of the dataset to share (required)")
	fs.StringVar(&users, "users", "", "Comma-separated list of usernames to share the dataset with (required)")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: hazo sharedb -creator <user> -mode <read|write> -dataset <name> -users <user1,user2,...>\n\n")
		fmt.Fprintf(fs.Output(), "Share a dataset with other users.\n\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return err
	}

	if creator == "" || mode == "" || datasetName == "" || users == "" {
		fs.Usage()
		return ErrInvalidArgs
	}
	if mode != "read" && mode != "write" {
		return fmt.Errorf("invalid mode %s: %w", mode, ErrInvalidMode)
	}

	sharedToUsers := strings.Split(users, ",")
	for i := range sharedToUsers {
		sharedToUsers[i] = strings.TrimSpace(sharedToUsers[i])
	}
	if len(sharedToUsers) == 0 {
		return ErrInvalidArgs
	}
	ctx := context.Background()
	db, queries, err := app.OpenDb()
	if err != nil {
		return fmt.Errorf("could not open users database: %w", err)
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := queries.WithTx(tx)
	dsRef := uuid.New().String()
	fmt.Printf("Creating dataset sharing '%s' ('%s') from user '%s' with mode '%s'\n", datasetName, dsRef, creator, mode)
	_, err = qtx.InsertDatasetSharing(ctx, appdb.InsertDatasetSharingParams{
		Ref:              dsRef,
		CreatorUserLogin: creator,
		CreationDate:     time.Now().Format(time.RFC3339),
		Name:             datasetName,
		Details:          sql.NullString{String: "Shared dataset", Valid: true},
	})
	if err != nil {
		return fmt.Errorf("could not insert dataset sharing '%s' ('%s', '%s'): %w", datasetName, dsRef, creator, err)
	}
	for _, user := range sharedToUsers {
		_, err = qtx.InsertDatasetSharingUser(ctx, appdb.InsertDatasetSharingUserParams{
			DatasetRef:          dsRef,
			DatasetCreatorLogin: creator,
			UserLogin:           user,
			Mode:                mode,
		})
		if err != nil {
			return fmt.Errorf("could not insert dataset sharing '%s' ('%s' of '%s') for user '%s': %w", datasetName, dsRef, creator, user, err)
		}
	}
	return tx.Commit()
}
