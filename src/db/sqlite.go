package db

import (
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Dataset struct {
	Name         string
	Path         string
	LastModified string
}

func ListDatasets() ([]Dataset, error) {
	files, err := filepath.Glob("./*.sq3")
	if err != nil {
		return nil, err
	}
	ds := make([]Dataset, len(files))
	for i, path := range files {
		info, err := os.Stat(path)
		if err != nil {
			return nil, err
		}
		info.ModTime()
		ds[i].Name = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		ds[i].Path = path
		ds[i].LastModified = info.ModTime().Format("2006-01-02 15:04:05")
	}
	return ds, nil
}

func ExecSqlite(dbPath string, code string) error {
	cmd := exec.Command("sqlite3", dbPath)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.WriteString(stdin, code)
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.Wait()
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
