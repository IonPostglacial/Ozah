package views

import _ "embed"

//go:embed taxon_form.html
var TaxonFormTemplate string

type TaxonFormData struct {
	Name        string
	NameV       string
	NameCN      string
	Author      string
	Website     string
	Description string
}
