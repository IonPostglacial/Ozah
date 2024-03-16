package main

import (
	"fmt"
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
	switch command := os.Args[1]; command {
	case "init":
		err = cmd.Init(os.Args[2:])
	case "import":
		err = cmd.ImportCsv(os.Args[2:])
	case "lsdoc":
		err = cmd.LsDoc(os.Args[2:])
	case "serve":
		err = cmd.Serve(os.Args[2:])
	default:
		err = fmt.Errorf("unknown command: '%s'", command)
	}
	if err != nil {
		log.Fatal(err)
	}
}
