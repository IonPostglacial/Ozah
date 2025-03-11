package storage

type Dataset struct {
	Name         string
	Path         string
	LastModified string
}

type PrivateDataset string

var InvalidPrivateDataset = PrivateDataset("\x00")
