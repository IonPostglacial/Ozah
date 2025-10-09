package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"nicolas.galipot.net/hazo/storage/app"
	"nicolas.galipot.net/hazo/storage/appdb"
)

const Cost = 11

func AddUser(args []string) error {
	fs := flag.NewFlagSet("adduser", flag.ExitOnError)

	var login, password, folderPath, capabilities string
	fs.StringVar(&login, "login", "", "Username for the new user (required)")
	fs.StringVar(&password, "password", "", "Password for the new user (will be hashed using bcrypt) (required)")
	fs.StringVar(&folderPath, "folder", "", "Private directory path for the user's data (required)")
	fs.StringVar(&capabilities, "capabilities", "", "Comma-separated list of capabilities to grant (e.g., 'user.manage,dataset.admin')")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: hazo adduser -login <username> -password <password> -folder <path> [-capabilities <list>]\n\n")
		fmt.Fprintf(fs.Output(), "Add a new user to the system.\n\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return err
	}

	if login == "" || password == "" || folderPath == "" {
		fs.Usage()
		return fmt.Errorf("all flags are required: -login, -password, -folder")
	}
	if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
		return fmt.Errorf("could not create directory '%s': %w", folderPath, err)
	}
	ctx := context.Background()
	_, queries, err := app.OpenDb()
	if err != nil {
		return fmt.Errorf("could not open users database: %w", err)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), Cost)
	if err != nil {
		return fmt.Errorf("could not hash password: %w", err)
	}
	_, err = queries.InsertCredentials(ctx, appdb.InsertCredentialsParams{
		Login:      login,
		Encryption: "bcrypt",
		Password:   string(hash),
	})
	if err != nil {
		return fmt.Errorf("could not insert credentials of user '%s': %w", login, err)
	}
	_, err = queries.InsertUserConfiguration(ctx, appdb.InsertUserConfigurationParams{
		Login:            login,
		PrivateDirectory: folderPath,
	})
	if err != nil {
		return fmt.Errorf("could not insert configuration of user '%s': %w", login, err)
	}
	if capabilities != "" {
		capList := splitCapabilities(capabilities)
		grantedDate := time.Now().Format(time.RFC3339)
		for _, cap := range capList {
			_, err = queries.GrantUserCapability(ctx, appdb.GrantUserCapabilityParams{
				UserLogin:      login,
				CapabilityName: cap,
				GrantedDate:    grantedDate,
				GrantedBy:      login,
			})
			if err != nil {
				return fmt.Errorf("could not grant capability '%s' to user '%s': %w", cap, login, err)
			}
		}
	}
	return nil
}

func splitCapabilities(capabilities string) []string {
	var result []string
	for _, cap := range strings.Split(capabilities, ",") {
		trimmed := strings.TrimSpace(cap)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
