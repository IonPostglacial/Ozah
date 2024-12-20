package summary

import (
	"context"
	"strings"

	"nicolas.galipot.net/hazo/db"
	"nicolas.galipot.net/hazo/db/storage"
)

type MultilangText struct {
	Name string
	Tr1  string
	Tr2  string
}

type State struct {
	Label MultilangText
	Color string
}

type Descriptor struct {
	Label  MultilangText
	States []State
}

type Section struct {
	Label       MultilangText
	Color       string
	Descriptors []Descriptor
}

type ViewModel struct {
	Sections []Section
}

func LoadForTaxon(ctx context.Context, queries *db.Queries, taxonRef string) (*ViewModel, error) {
	sd, err := queries.GetSummaryDescriptors(ctx, taxonRef)
	if err != nil {
		return nil, err
	}
	summary := &ViewModel{}
	statesByPath := make(map[string][]State)
	for _, state := range sd {
		statesByPath[state.Path] = append(statesByPath[state.Path], State{
			Label: MultilangText{
				Name: state.Name,
				Tr1:  state.NameTr1.String,
				Tr2:  state.NameTr2.String,
			},
			Color: state.Color.String,
		})
	}
	characterRefs := make([]string, 0, len(statesByPath))
	for path := range statesByPath {
		parentId := path[strings.LastIndex(path, ".")+1:]
		characterRefs = append(characterRefs, parentId)
	}
	characters, err := queries.GetCatCharactersNameTr2(ctx, storage.GetCatCharactersNameTr2Params{
		Lang1: "EN", Lang2: "CN", Refs: characterRefs,
	})
	if err != nil {
		return nil, err
	}
	descriptionsBySection := make(map[string][]Descriptor)
	for _, ch := range characters {
		fullPath := db.FullPath(ch.Path, ch.Ref)
		states := statesByPath[fullPath]
		path := strings.Split(ch.Path, ".")
		section := "c0"
		if len(path) > 1 {
			section = path[1]
		}
		descriptionsBySection[section] = append(descriptionsBySection[section], Descriptor{
			Label: MultilangText{
				Name: ch.Name,
				Tr1:  ch.NameTr1.String,
				Tr2:  ch.NameTr2.String,
			},
			States: states,
		})
	}
	sectionRefs := make([]string, 0, len(descriptionsBySection))
	for ref := range descriptionsBySection {
		sectionRefs = append(sectionRefs, ref)
	}
	sections, err := queries.GetCatCharactersNameTr2(ctx, storage.GetCatCharactersNameTr2Params{
		Lang1: "EN", Lang2: "CN", Refs: sectionRefs,
	})
	if err != nil {
		return nil, err
	}
	for _, sec := range sections {
		descriptors := descriptionsBySection[sec.Ref]
		summary.Sections = append(summary.Sections, Section{
			Label: MultilangText{
				Name: sec.Name,
				Tr1:  sec.NameTr1.String,
				Tr2:  sec.NameTr2.String,
			},
			Color:       sec.Color.String,
			Descriptors: descriptors,
		})
	}
	return summary, nil
}
