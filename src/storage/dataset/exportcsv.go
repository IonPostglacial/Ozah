package dataset

import (
	"archive/zip"
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

func ExportCsv(dsName string, queries *Queries, w io.Writer) error {
	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()

	ctx := context.Background()
	db := queries.db

	for _, tableName := range tableNames {
		if err := exportTable(ctx, db, zipWriter, tableName); err != nil {
			return fmt.Errorf("failed to export table '%s': %w", tableName, err)
		}
	}

	return nil
}

func exportTable(ctx context.Context, db *sql.DB, zipWriter *zip.Writer, tableName string) error {
	columns, err := getTableColumns(ctx, db, tableName)
	if err != nil {
		return fmt.Errorf("could not get columns for table '%s': %w", tableName, err)
	}

	query := fmt.Sprintf("SELECT * FROM %s", tableName)
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil
	}
	defer rows.Close()

	fileName := fmt.Sprintf("%s.csv", tableName)
	csvFile, err := zipWriter.Create(fileName)
	if err != nil {
		return fmt.Errorf("could not create CSV file '%s' in ZIP: %w", fileName, err)
	}

	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()

	header := make([]string, len(columns))
	for i, col := range columns {
		header[i] = fmt.Sprintf("%s:%s", col.Name, col.Type)
	}
	if err := csvWriter.Write(header); err != nil {
		return fmt.Errorf("could not write CSV header: %w", err)
	}

	rowValues := make([]any, len(columns))
	rowPointers := make([]any, len(columns))
	for i := range rowValues {
		rowPointers[i] = &rowValues[i]
	}

	for rows.Next() {
		if err := rows.Scan(rowPointers...); err != nil {
			return fmt.Errorf("could not scan row: %w", err)
		}

		rowStrings := make([]string, len(columns))
		for i, val := range rowValues {
			rowStrings[i] = formatCsvValue(val)
		}

		if err := csvWriter.Write(rowStrings); err != nil {
			return fmt.Errorf("could not write CSV row: %w", err)
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating rows: %w", err)
	}

	return nil
}

type tableColumn struct {
	Name string
	Type string
}

func getTableColumns(ctx context.Context, db *sql.DB, tableName string) ([]tableColumn, error) {
	query := fmt.Sprintf("PRAGMA table_info(%s)", tableName)
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []tableColumn
	for rows.Next() {
		var cid int
		var name, colType string
		var notNull int
		var dfltValue sql.NullString
		var pk int

		if err := rows.Scan(&cid, &name, &colType, &notNull, &dfltValue, &pk); err != nil {
			return nil, err
		}

		columns = append(columns, tableColumn{
			Name: strings.ToLower(name),
			Type: strings.ToLower(colType),
		})
	}

	return columns, rows.Err()
}

func formatCsvValue(val any) string {
	if val == nil {
		return ""
	}

	switch v := val.(type) {
	case string:
		return v
	case int64:
		return fmt.Sprintf("%d", v)
	case float64:
		return fmt.Sprintf("%g", v)
	case bool:
		if v {
			return "true"
		}
		return "false"
	case []byte:
		return string(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}
