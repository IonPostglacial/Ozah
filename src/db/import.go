package db

import (
	_ "embed"
	"encoding/csv"
	"fmt"
	"html/template"
	"io"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
)

type tmplData struct {
	DirectoryPath string
}

//go:embed import_template.txt
var tmpl string

var tableNames = []string{
	"lang", "unit",
	"document", "document_translation", "document_attachment", "book",
	"measurement_character", "categorical_character", "state",
	"periodic_character", "geographical_place", "geographical_map", "geographical_character",
	"descriptor_visibility_inapplicable", "descriptor_visibility_requirement",
	"taxon", "taxon_book_info", "taxon_measurement", "taxon_description", "taxon_specimen_location",
}

func ImportCsv(csvPath string, ds PrivateDataset) error {
	for _, tableName := range tableNames {
		err := importFile(csvPath, ds, tableName)
		if err != nil {
			return fmt.Errorf("error importing CSV '%s': %w", csvPath, err)
		}
	}
	tmpl, err := template.New("import_template").Parse(tmpl)
	if err != nil {
		return err
	}
	var buf strings.Builder
	err = tmpl.Execute(&buf, tmplData{DirectoryPath: csvPath})
	if err != nil {
		return err
	}
	return ExecSqlite(string(ds), buf.String())
}

func importFile(csvPath string, ds PrivateDataset, tableName string) error {
	fileName := fmt.Sprintf("%s.csv", tableName)
	filePath := path.Join(csvPath, fileName)
	content, err := parseCsvFile(filePath)
	if err != nil {
		return fmt.Errorf("error importing file '%s': %w", filePath, err)
	}
	fmt.Printf("content: %#v\n", content)
	return nil
}

type columnParser = func(text string) (any, error)

type column struct {
	Name  string
	Parse columnParser
}

type csvContent struct {
	columns []column
	rows    [][]any
}

var ErrInvalidColor = fmt.Errorf("invalid color format, expected hexadecimal rgb(a)")
var ErrInvalidUrl = fmt.Errorf("invalid URL format")
var ErrUnknownColumnType = fmt.Errorf("unknown column type")

func parseString(text string) (any, error) {
	return text, nil
}

func parseInt64(text string) (any, error) {
	if text == "" {
		return nil, nil
	}
	return strconv.ParseInt(text, 10, 64)
}

func parseFloat64(text string) (any, error) {
	if text == "" {
		return nil, nil
	}
	return strconv.ParseFloat(text, 64)
}

func parseColor(text string) (any, error) {
	if text == "" {
		return nil, nil
	}
	if len(text) < 3 || len(text) > 9 || text[0] != '#' {
		return nil, fmt.Errorf("unexpected text '%s': %w", text, ErrInvalidColor)
	}
	_, err := strconv.ParseUint(text[1:], 16, 64)
	if err != nil {
		return nil, fmt.Errorf("unexpected text '%s': %w", text, ErrInvalidColor)
	}
	return text, nil
}

func parseUrl(text string) (any, error) {
	if text == "" {
		return nil, nil
	}
	_, err := url.Parse(text)
	if err != nil {
		return nil, fmt.Errorf("unexpected text '%s': %w", text, ErrInvalidUrl)
	}
	return text, nil
}

func parseHeader(r *csv.Reader) ([]column, error) {
	header, err := r.Read()
	if err != nil {
		return nil, err
	}
	columns := make([]column, len(header))
	for i, colName := range header {
		nameType := strings.SplitN(colName, ":", 2)
		if len(nameType) <= 1 {
			columns[i].Name = colName
			columns[i].Parse = parseString
			continue
		}
		columns[i].Name = nameType[0]
		switch nameType[1] {
		case "i64":
			columns[i].Parse = parseInt64
		case "f64":
			columns[i].Parse = parseFloat64
		case "hexcolor":
			columns[i].Parse = parseColor
		case "url":
			columns[i].Parse = parseUrl
		case "ltree":
			columns[i].Parse = parseString
		case "path":
			columns[i].Parse = parseString
		default:
			return nil, fmt.Errorf("invalid CSV header column type '%s': %w", nameType[1], ErrUnknownColumnType)
		}
	}
	return columns, nil
}

func parseCsvFile(fileName string) (*csvContent, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("error opening CSV file to import: %w", err)
	}
	defer f.Close()
	r := csv.NewReader(f)
	r.TrimLeadingSpace = true
	columns, err := parseHeader(r)
	if err != nil {
		return nil, fmt.Errorf("error reading CSV header: %w", err)
	}
	rows := make([][]any, 0)
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading CSV record: %w", err)
		}
		row := make([]any, len(columns))
		for i, text := range rec {
			row[i], err = columns[i].Parse(text)
			if err != nil {
				return nil, fmt.Errorf("error parsing value '%s' for column '%s': %w", text, columns[i].Name, err)
			}
		}
		rows = append(rows, row)
	}
	return &csvContent{columns, rows}, nil
}
