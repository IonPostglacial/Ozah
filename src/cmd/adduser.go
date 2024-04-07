package cmd

import (
	"context"
	"os"

	"golang.org/x/crypto/bcrypt"
	"nicolas.galipot.net/hazo/db"
	"nicolas.galipot.net/hazo/db/commonstorage"
)

const Cost = 11

func AddUser(args []string) error {
	login := args[0]
	password := args[1]
	folderPath := args[2]
	if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
		return err
	}
	ctx := context.Background()
	_, queries, err := db.OpenCommon()
	if err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), Cost)
	if err != nil {
		return err
	}
	_, err = queries.InsertCredentials(ctx, commonstorage.InsertCredentialsParams{
		Login:      login,
		Encryption: "bcrypt",
		Password:   string(hash),
	})
	if err != nil {
		return err
	}
	_, err = queries.InsertUserConfiguration(ctx, commonstorage.InsertUserConfigurationParams{
		Login:            login,
		PrivateDirectory: folderPath,
	})
	if err != nil {
		return err
	}
	return nil
}
