package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"nicolas.galipot.net/hazo/cmd"
)

const usage = `Hazo - A botanical taxonomy management system

Usage:
  hazo [command] [options]

Available Commands:
  setup       Initialize the application database with default configuration
  init        Create a new dataset database
  adduser     Add a new user to the system
  sharedb     Share a dataset with other users
  serve       Start the web server
  importcsv   Import data from a CSV file into a dataset
  importjson  Import data from a JSON file into a dataset
  exportjson  Export a dataset to JSON format
  lsdoc       List documents in a dataset
  help        Show this help message

Use "hazo [command] -h" for more information about a command.

If no command is provided, the server will start on :8080 with debug mode enabled.
`

func printUsage() {
	fmt.Fprint(os.Stderr, usage)
}

func main() {
	var err error
	if len(os.Args) < 2 {
		err = cmd.Serve([]string{":8080"})
		log.Fatal(err)
	}
	switch command := os.Args[1]; command {
	case "help", "-h", "--help":
		printUsage()
		return
	case "setup":
		err = cmd.Setup(os.Args[2:])
	case "init":
		err = cmd.Init(os.Args[2:])
	case "adduser":
		err = cmd.AddUser(os.Args[2:])
	case "sharedb":
		err = cmd.Sharedb(os.Args[2:])
	case "importcsv":
		err = cmd.ImportCsv(os.Args[2:])
	case "importjson":
		err = cmd.ImportJson(os.Args[2:])
	case "exportjson":
		err = cmd.ExportJson(os.Args[2:])
	case "lsdoc":
		err = cmd.LsDoc(os.Args[2:])
	case "serve":
		err = cmd.Serve(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "Error: unknown command '%s'\n\n", command)
		printUsage()
		os.Exit(1)
	}
	if err != nil {
		log.Fatal(err)
	}
}
