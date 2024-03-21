package db

import (
	_ "embed"
)

//go:embed schema.sql
var Schema string

//go:embed index.sql
var Index string

//go:embed common_schema.sql
var CommonSchema string
