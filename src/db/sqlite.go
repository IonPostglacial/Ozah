package db

import (
	"fmt"
	"io"
	"os/exec"
)

type Dataset struct {
	Name         string
	Path         string
	LastModified string
}

func ExecSqlite(dbPath string, code string) error {
	cmd := exec.Command("sqlite3", dbPath)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("piping stdin to sqlite3 failed: %w", err)
	}
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("launching sqlite3 process failed: %w", err)
	}
	_, err = io.WriteString(stdin, code)
	if err != nil {
		return fmt.Errorf("writing the query to stdin failed: %w", err)
	}
	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("waiting for sqlite3 to complete failed: %w", err)
	}
	return nil
}
