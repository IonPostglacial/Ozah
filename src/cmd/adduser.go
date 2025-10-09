package cmd

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"nicolas.galipot.net/hazo/user"
)

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

	ctx := context.Background()
	var capList []string
	if capabilities != "" {
		capList = splitCapabilities(capabilities)
	}

	return user.Create(ctx, user.CreateUserParams{
		Login:            login,
		Password:         password,
		PrivateDirectory: folderPath,
		Capabilities:     capList,
		GrantedBy:        login,
	})
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
