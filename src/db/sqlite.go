package db

import (
	"io"
	"log"
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
