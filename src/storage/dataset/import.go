package dataset

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

func ImportJsonDataset(dbPath Private, jsonData []byte) error {
	if err := Create(dbPath); err != nil {
		return fmt.Errorf("could not create dataset: %w", err)
	}

	if err := ImportJson(jsonData, dbPath); err != nil {
		return fmt.Errorf("could not import JSON data: %w", err)
	}

	return nil
}

func ImportCsvDataset(dbPath Private, zipData []byte) error {
	dir, err := os.MkdirTemp("", "csv-import-")
	if err != nil {
		return fmt.Errorf("could not create temporary directory: %w", err)
	}
	defer os.RemoveAll(dir)

	zipReader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return fmt.Errorf("could not read ZIP file: %w", err)
	}

	for _, f := range zipReader.File {
		if strings.HasSuffix(f.Name, "/") {
			continue
		}

		content, err := f.Open()
		if err != nil {
			return fmt.Errorf("could not read file '%s' from ZIP: %w", f.Name, err)
		}
		defer content.Close()

		filePath := path.Join(dir, path.Base(f.Name))
		if err := os.MkdirAll(path.Dir(filePath), 0770); err != nil {
			return fmt.Errorf("could not create directory for '%s': %w", f.Name, err)
		}

		outFile, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("could not create file '%s': %w", filePath, err)
		}
		defer outFile.Close()

		if _, err = io.Copy(outFile, content); err != nil {
			return fmt.Errorf("could not write file '%s': %w", filePath, err)
		}
	}

	if err := Create(dbPath); err != nil {
		return fmt.Errorf("could not create dataset: %w", err)
	}

	if err := ImportCsv(dir, dbPath); err != nil {
		return fmt.Errorf("could not import CSV data: %w", err)
	}

	return nil
}
