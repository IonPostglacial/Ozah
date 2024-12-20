package identification

import (
	"nicolas.galipot.net/hazo/server/components/popover"
)

type TaxonViewModel struct {
	Ref  string
	Name string
	Url  string
}

type SelectedState struct {
	ParentRef  string
	ParentName string
	Ref        string
	Name       string
	Url        string
}

type State struct {
	Ref  string
	Name string
	Url  string
}

type CharacterViewModel struct {
	Ref    string
	Name   string
	States []State
}

type MeasurementViewModel struct {
	Ref      string
	Name     string
	UnitRef  string
	UnsetUrl string
	HasValue bool
	Value    float64
}

type ViewModel struct {
	PageTitle             string
	AvailableDatasets     *popover.ViewModel
	ViewMenuState         *popover.ViewModel
	Taxa                  []TaxonViewModel
	Characters            []CharacterViewModel
	MeasurementCharacters []MeasurementViewModel
	SelectedStates        []SelectedState
}
