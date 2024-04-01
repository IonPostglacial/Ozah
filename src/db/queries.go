package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"nicolas.galipot.net/hazo/db/storage"
)

const identifyTaxonPrelude = `select doc.Ref, doc.Name from Document doc
where doc.Ref in (`

const identifyTaxonWithMeasurements = `
	with specimen(Character_Ref, Measurement) as (values %s) 
	select m.Taxon_Ref from specimen 
	left join Taxon_Measurement m on m.Character_Ref = specimen.Character_Ref 
	where 
		(Minimum is null or specimen.Measurement >= Minimum) and 
		(Maximum is null or specimen.Measurement <= Maximum)
`

const identifyTaxonWithStates = `
	select Taxon_Ref from Taxon_Description 
    where Description_Ref in (%s) 
    group by Taxon_Ref 
    having Count(Taxon_Ref) = ?
`

const identifyTaxonCoda = `
) order by doc.Name asc;`

type SpecimenMeasurement struct {
	CharacterRef string
	Value        float64
}

type TaxonIdentificationParams struct {
	Measurements []SpecimenMeasurement
	StateRefs    []string
}

type IdentifiedTaxon struct {
	Ref  string
	Name string
}

type Queries struct {
	*storage.Queries
	db *sql.DB
}

func (q *Queries) IdentifyTaxa(ctx context.Context, arg TaxonIdentificationParams) ([]IdentifiedTaxon, error) {
	var measurementValues, states, query strings.Builder
	var measurementSep, stateSep string
	hasMeasurement := len(arg.Measurements) > 0
	hasStates := len(arg.StateRefs) > 0
	args := make([]any, 0, len(arg.Measurements)+len(arg.StateRefs)+1)
	for _, measurement := range arg.Measurements {
		measurementValues.WriteString(measurementSep)
		measurementValues.WriteString("(?,?)")
		args = append(args, measurement.CharacterRef, measurement.Value)
		measurementSep = ","
	}
	for _, stateRef := range arg.StateRefs {
		states.WriteString(stateSep)
		states.WriteString("?")
		args = append(args, stateRef)
		stateSep = ","
	}
	if hasStates {
		args = append(args, len(arg.StateRefs))
	}
	query.WriteString(identifyTaxonPrelude)
	switch {
	case hasMeasurement && hasStates:
		fmt.Fprintf(&query, identifyTaxonWithMeasurements, measurementValues.String())
		query.WriteString(" intersect ")
		fmt.Fprintf(&query, identifyTaxonWithStates, states.String())
	case hasMeasurement:
		fmt.Fprintf(&query, identifyTaxonWithMeasurements, measurementValues.String())
	case hasStates:
		fmt.Fprintf(&query, identifyTaxonWithStates, states.String())
	default:
		return []IdentifiedTaxon{}, nil
	}
	query.WriteString(identifyTaxonCoda)
	rows, err := q.db.QueryContext(ctx, query.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []IdentifiedTaxon
	for rows.Next() {
		var i IdentifiedTaxon
		if err := rows.Scan(
			&i.Ref,
			&i.Name,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
