package db

import (
	_ "embed"
)

//go:embed common_schema.sql
var CommonSchema string
