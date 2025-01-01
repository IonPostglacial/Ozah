package db

import (
	_ "embed"
)

//go:embed schema.sql
var schema string

//go:embed index.sql
var index string
