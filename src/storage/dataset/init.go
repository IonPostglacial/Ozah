package dataset

import (
	_ "embed"
)

//go:embed schema.sql
var Schema string

//go:embed index.sql
var Index string
