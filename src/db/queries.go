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

const getDocumentsHierarchyQueryPrelude = `select doc.ref, doc.path, doc.doc_order, doc.name, doc.details`
const getDocumentsHierarchyQueryIntro = ` from Document doc`
const getDocumentsHierarchyTrQuery = `
left join Document_Translation tr%[1]d on doc.Ref = tr%[1]d.Document_Ref and tr%[1]d.Lang_Ref = ?%[2]d`
const getDocumentsHierarchyQueryCoda = `
where (doc.Path >= ?1 and doc.Path < (?1 || '.~'))
order by (Path || '.' || Ref) asc;`

type Document struct {
	Ref      string
	Path     string
	DocOrder int64
	Name     string
	NameTr   []sql.NullString
	Details  string
}

func (q *Queries) GetDocumentHierarchy(ctx context.Context, documentPath string, langs []string) ([]Document, error) {
	var query, queryTr strings.Builder
	query.WriteString(getDocumentsHierarchyQueryPrelude)
	args := make([]any, len(langs)+1)
	args[0] = documentPath
	for i, lang := range langs {
		args[i+1] = lang
		query.WriteString(fmt.Sprintf(", tr%[1]d.name name_tr%[1]d", i+1))
		queryTr.WriteString(fmt.Sprintf(getDocumentsHierarchyTrQuery, i+1, i+2))
	}
	query.WriteString(getDocumentsHierarchyQueryIntro)
	query.WriteString(queryTr.String())
	query.WriteString(getDocumentsHierarchyQueryCoda)

	rows, err := q.db.QueryContext(ctx, query.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Document
	for rows.Next() {
		doc := Document{
			NameTr: make([]sql.NullString, len(langs)),
		}
		etc := make([]any, 5+len(langs))
		etc[0] = &doc.Ref
		etc[1] = &doc.Path
		etc[2] = &doc.DocOrder
		etc[3] = &doc.Name
		etc[4] = &doc.Details
		for i := range doc.NameTr {
			etc[5+i] = &doc.NameTr[i]
		}
		if err := rows.Scan(etc...); err != nil {
			return nil, err
		}
		items = append(items, doc)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
