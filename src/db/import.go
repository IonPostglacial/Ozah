package db

import (
	_ "embed"
	"strings"
	"text/template"
)

type tmplData struct {
	DirectoryPath string
}

//go:embed import_template.txt
var tmpl string

func ImportCsv(csvPath string, to string) error {
	tmpl, err := template.New("import_template").Parse(tmpl)
	if err != nil {
		return err
	}
	var buf strings.Builder
	err = tmpl.Execute(&buf, tmplData{DirectoryPath: csvPath})
	if err != nil {
		return err
	}
	return ExecSqlite(to, buf.String())
}
