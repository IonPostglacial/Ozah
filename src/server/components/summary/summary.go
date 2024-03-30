package summary

import (
	"context"
	"fmt"
	"strings"

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
	Color  string
	States []State
}

type Model struct {
	Descriptors []Descriptor
}

func LoadForTaxon(ctx context.Context, queries *storage.Queries, taxonRef string) (*Model, error) {
	sd, err := queries.GetSummaryDescriptors(ctx, taxonRef)
	if err != nil {
		return nil, err
	}
	summary := &Model{}
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
	characterIds := make([]string, 0, len(statesByPath))
	for path := range statesByPath {
		parentId := path[strings.LastIndex(path, ".")+1:]
		characterIds = append(characterIds, parentId)
	}
	characters, err := queries.GetCatCharactersNameTr2(ctx, storage.GetCatCharactersNameTr2Params{
		Lang1: "EN", Lang2: "CN", Refs: characterIds,
	})
	if err != nil {
		return nil, err
	}
	for _, ch := range characters {
		fullPath := fmt.Sprintf("%s.%s", ch.Path, ch.Ref)
		states := statesByPath[fullPath]
		summary.Descriptors = append(summary.Descriptors, Descriptor{
			Label: MultilangText{
				Name: ch.Name,
				Tr1:  ch.NameTr1.String,
				Tr2:  ch.NameTr2.String,
			},
			Color:  ch.Color.String,
			States: states,
		})
	}
	return summary, nil
}
