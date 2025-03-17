package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"nicolas.galipot.net/hazo/cmd"
	"nicolas.galipot.net/hazo/server/common"
)

func main() {
	var err error
	if len(os.Args) < 2 {
		err = cmd.Serve([]string{":8080"}, &common.ServerConfig{Debug: true})
		log.Fatal(err)
	}
	switch command := os.Args[1]; command {
	case "setup":
		err = cmd.Setup(os.Args[2:])
	case "init":
		err = cmd.Init(os.Args[2:])
	case "adduser":
		err = cmd.AddUser(os.Args[2:])
	case "importcsv":
		err = cmd.ImportCsv(os.Args[2:])
	case "importjson":
		err = cmd.ImportJson(os.Args[2:])
	case "lsdoc":
		err = cmd.LsDoc(os.Args[2:])
	case "serve":
		err = cmd.Serve(os.Args[2:], &common.ServerConfig{})
	default:
		err = fmt.Errorf("unknown command: '%s'", command)
	}
	if err != nil {
		log.Fatal(err)
	}
}
