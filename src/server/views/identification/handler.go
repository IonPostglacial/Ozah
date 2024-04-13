package identification

import (
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"maps"

	"nicolas.galipot.net/hazo/db"
	"nicolas.galipot.net/hazo/server/common"
	"nicolas.galipot.net/hazo/server/components"
	"nicolas.galipot.net/hazo/server/components/popover"
	"nicolas.galipot.net/hazo/server/link"
	"nicolas.galipot.net/hazo/server/views"
)

//go:embed identification.html
var identificationPage string

type Model struct {
	PageTitle             string
	AvailableDatasets     *popover.State
	ViewMenuState         *popover.State
	IdentifiedTaxa        []IdentifiedTaxon
	Characters            []Character
	MeasurementCharacters []Measurement
	SelectedStates        []SelectedState
}

type IdentifiedTaxon struct {
	db.IdentifiedTaxon
	Url string
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

type Character struct {
	Ref    string
	Name   string
	States []State
}

type Measurement struct {
	Ref      string
	Name     string
	UnitRef  string
	UnsetUrl string
	HasValue bool
	Value    float64
}

func LinkToIdentification(dsName string, stateRefs []string, measures map[string]db.SpecimenMeasurement) string {
	linkUrl := strings.Builder{}
	urlQuery := url.Values{}
	for _, ref := range stateRefs {
		urlQuery.Add("s", ref)
	}
	for _, measure := range measures {
		urlQuery.Add(fmt.Sprintf("m-%s", measure.CharacterRef), fmt.Sprintf("%f", measure.Value))
	}
	linkUrl.WriteString(link.ToIdentify(dsName))
	linkUrl.WriteRune('?')
	linkUrl.WriteString(urlQuery.Encode())
	return linkUrl.String()
}

func Handler(w http.ResponseWriter, r *http.Request, cc *common.Context) error {
	dsName := r.PathValue("dsName")
	ctx := context.Background()
	cc.Template = components.NewTemplate()
	ds, err := cc.User.GetDataset(dsName)
	if err != nil {
		return err
	}
	queries, err := db.Open(ds)
	if err != nil {
		return err
	}
	queryParams := r.URL.Query()
	stateRefs := queryParams["s"]
	measurements := make(map[string]db.SpecimenMeasurement, len(queryParams))
	for k, values := range queryParams {
		if strings.HasPrefix(k, "m-") {
			for _, v := range values {
				value, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return fmt.Errorf("wrong measurement query '%s': %w", k, err)
				}
				measurements[k[2:]] = db.SpecimenMeasurement{
					CharacterRef: k[2:],
					Value:        value,
				}
			}
		}
	}
	statesInfo, err := queries.GetDocumentsAndParentsNames(ctx, stateRefs)
	if err != nil {
		return err
	}
	selectedStates := make([]SelectedState, len(statesInfo))
	for i, state := range statesInfo {
		refs := slices.Clone(stateRefs)
		refs = slices.DeleteFunc(refs, func(s string) bool { return s == state.Ref })
		url := LinkToIdentification(dsName, refs, measurements)
		selectedStates[i] = SelectedState{
			ParentRef: state.ParentRef.String, ParentName: state.ParentName.String,
			Ref: state.Ref, Name: state.Name, Url: url,
		}
	}
	ms := make([]db.SpecimenMeasurement, 0, len(measurements))
	for _, v := range measurements {
		ms = append(ms, v)
	}
	taxa, err := queries.IdentifyTaxa(ctx, db.TaxonIdentificationParams{
		Measurements: ms,
		StateRefs:    stateRefs,
	})
	if err != nil {
		return fmt.Errorf("error executing identification query: %w", err)
	}
	identifiedTaxa := make([]IdentifiedTaxon, len(taxa))
	for i, taxon := range taxa {
		identifiedTaxa[i] = IdentifiedTaxon{taxon, link.ToTaxon(dsName, taxon.Ref)}
	}
	datasets, err := views.NewDatasetMenuState(cc, dsName)
	if err != nil {
		return err
	}
	distinctiveCharacters, err := queries.DistinctiveCharacters(ctx)
	if err != nil {
		return err
	}
	chars := make([]Character, 0)
	var lastChar *Character
	for _, ch := range distinctiveCharacters {
		if lastChar == nil || ch.Ref != lastChar.Ref {
			chars = append(chars, Character{
				Ref:  ch.Ref,
				Name: ch.Name,
			})
			lastChar = &chars[len(chars)-1]
		}
		refs := slices.Clone(stateRefs)
		refs = append(refs, ch.StateRef)
		url := LinkToIdentification(dsName, refs, measurements)
		state := State{Ref: ch.StateRef, Name: ch.StateName, Url: url}
		lastChar.States = append(lastChar.States, state)
	}
	mcs, err := queries.GetMeasurementCharacters(ctx)
	if err != nil {
		return err
	}
	measurementChars := make([]Measurement, len(mcs))
	for i, mc := range mcs {
		value := queryParams.Get(fmt.Sprintf("m-%s", mc.Ref))
		val, err := strconv.ParseFloat(value, 64)
		ms := maps.Clone(measurements)
		delete(ms, mc.Ref)
		measurementChars[i] = Measurement{
			Ref: mc.Ref, Name: mc.Name, UnitRef: mc.UnitRef.String,
			UnsetUrl: LinkToIdentification(dsName, stateRefs, ms),
			HasValue: err == nil,
			Value:    val,
		}
	}
	template.Must(cc.Template.Parse(identificationPage))
	err = cc.Template.Execute(w, Model{
		PageTitle:             "Identification",
		AvailableDatasets:     datasets,
		ViewMenuState:         views.NewMenuState("Identify", dsName),
		IdentifiedTaxa:        identifiedTaxa,
		Characters:            chars,
		MeasurementCharacters: measurementChars,
		SelectedStates:        selectedStates,
	})
	if err != nil {
		return err
	}
	return nil
}
