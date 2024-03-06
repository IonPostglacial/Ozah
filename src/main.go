package main

import (
	"context"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"nicolas.galipot.net/hazo/cmd"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Pass command 'init' or 'run'.")
	}
	var err error
	ctx := context.Background()
	switch os.Args[1] {
	case "init":
		err = cmd.Init(os.Args[2])
	case "import":
		err = cmd.ImportCsv(os.Args[2], os.Args[3])
	case "lsdoc":
		err = cmd.LsDoc(ctx, os.Args[2], os.Args[3])
	case "serve":
		err = cmd.Serve(os.Args[2])
	default:
		log.Fatal("unknown command:", os.Args[1])
	}
	if err != nil {
		log.Fatal(err)
	}
}
