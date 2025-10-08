package dataset

type T struct {
	Name         string
	Path         string
	LastModified string
}

type Shared struct {
	T
	Creator string
	Mode    string // "read" or "write"
}

type Private string

var InvalidPrivate = Private("\x00")
