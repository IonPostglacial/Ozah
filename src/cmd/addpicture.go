package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"nicolas.galipot.net/hazo/storage/picture"
	"nicolas.galipot.net/hazo/user"
)

func AddPicture(args []string) error {
	fs := flag.NewFlagSet("addpicture", flag.ExitOnError)

	var datasetName, docRef, filePath, login string
	var attachmentIndex int
	fs.StringVar(&login, "login", "", "User login (required)")
	fs.StringVar(&datasetName, "dataset", "", "Name of the dataset (required)")
	fs.StringVar(&docRef, "ref", "", "Document reference (taxon/character/state ref) (required)")
	fs.StringVar(&filePath, "file", "", "Path to the picture file (required)")
	fs.IntVar(&attachmentIndex, "index", 0, "Attachment index (default: 0)")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: hazo addpicture -login <username> -dataset <name> -ref <id> -file <path> [-index <n>]\n\n")
		fmt.Fprintf(fs.Output(), "Add a picture to a document in a dataset with automatic thumbnail generation.\n\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return err
	}

	if login == "" || datasetName == "" || docRef == "" || filePath == "" {
		fs.Usage()
		return fmt.Errorf("required flags: -login, -dataset, -ref, -file")
	}

	u, err := user.Register(login)
	if err != nil {
		return fmt.Errorf("could not register user: %w", err)
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", filePath)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	ctx := context.Background()
	result, err := picture.UploadPicture(ctx, u, datasetName, docRef, file, filepath.Base(filePath), attachmentIndex)
	if err != nil {
		return fmt.Errorf("could not upload picture: %w", err)
	}

	fmt.Printf("Successfully added picture to document '%s' at index %d\n", docRef, result.AttachmentIndex)
	fmt.Printf("  Original: %s\n", result.OriginalPath)
	fmt.Printf("  Small: %s\n", result.SmallPath)
	fmt.Printf("  Medium: %s\n", result.MediumPath)
	fmt.Printf("  Big: %s\n", result.BigPath)

	return nil
}
