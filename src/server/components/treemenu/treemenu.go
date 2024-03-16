package treemenu

import _ "embed"

//go:embed entry.html
var EntryTemplate string

//go:embed treemenu.html
var Template string

type Item struct {
	Id       string
	Url      string
	FullPath string
	Order    int
	Name     string
	NameV    string
	NameCN   string
	Children []*Item
}
