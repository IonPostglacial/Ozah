package storage

type Dataset struct {
	Name         string
	Path         string
	LastModified string
}

type SharedDataset struct {
	Dataset
	Creator string
	Mode    string // "read" or "write"
}

type PrivateDataset string

var InvalidPrivateDataset = PrivateDataset("\x00")
