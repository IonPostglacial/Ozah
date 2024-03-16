package cmd

import (
	_ "embed"
	"strings"
	"text/template"

	"nicolas.galipot.net/hazo/db"
)

type tmplData struct {
	DirectoryPath string
}

//go:embed import_template.txt
var tmpl string

func ImportCsv(args []string) error {
	csvPath := args[0]
	to := args[1]
	tmpl, err := template.New("import_template").Parse(tmpl)
	if err != nil {
		return err
	}
	var buf strings.Builder
	err = tmpl.Execute(&buf, tmplData{DirectoryPath: csvPath})
	if err != nil {
		return err
	}
	db.ExecSqlite(to, buf.String())
	return nil
}
