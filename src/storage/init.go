package storage

import (
	_ "embed"
)

//go:embed dataset_schema.sql
var DatasetSchema string

//go:embed index.sql
var index string

//go:embed app_schema.sql
var AppSchema string
