package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"nicolas.galipot.net/hazo/storage/dsdb"
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
	*dsdb.Queries
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

func reserveBuffer(buf []byte, appendSize int) []byte {
	newSize := len(buf) + appendSize
	if cap(buf) < newSize {
		newBuf := make([]byte, len(buf)*2+appendSize)
		copy(newBuf, buf)
		buf = newBuf
	}
	return buf[:newSize]
}

func escapeBytesBackslash(buf []byte, v []byte) []byte {
	pos := len(buf)
	buf = reserveBuffer(buf, len(v)*2)

	for _, c := range v {
		switch c {
		case '\x00':
			buf[pos] = '\\'
			buf[pos+1] = '0'
			pos += 2
		case '\n':
			buf[pos] = '\\'
			buf[pos+1] = 'n'
			pos += 2
		case '\r':
			buf[pos] = '\\'
			buf[pos+1] = 'r'
			pos += 2
		case '\x1a':
			buf[pos] = '\\'
			buf[pos+1] = 'Z'
			pos += 2
		case '\'':
			buf[pos] = '\\'
			buf[pos+1] = '\''
			pos += 2
		case '"':
			buf[pos] = '\\'
			buf[pos+1] = '"'
			pos += 2
		case '\\':
			buf[pos] = '\\'
			buf[pos+1] = '\\'
			pos += 2
		default:
			buf[pos] = c
			pos++
		}
	}

	return buf[:pos]
}

func escapeStringBackslash(buf []byte, v string) []byte {
	return escapeBytesBackslash(buf, []byte(v))
}

func EscapeString(s string) string {
	buf := make([]byte, 0, len(s))
	return string(escapeStringBackslash(buf, s))
}

const getDocumentsHierarchyQueryPrelude = `select doc.ref, doc.path, doc.doc_order, doc.name, doc.details`
const getDocumentsHierarchyQueryIntro = ` from Document doc`
const getDocumentsHierarchyTrQuery = `
left join Document_Translation tr%[1]d on doc.Ref = tr%[1]d.Document_Ref and tr%[1]d.Lang_Ref = ?%[2]d`
const getDocumentsChildrenTrQuery = `
left join Document_Translation ctr%[1]d on children.Ref = ctr%[1]d.Document_Ref and ctr%[1]d.Lang_Ref = ?%[2]d`
const getDocumentsHierarchyQueryCoda = `
where (doc.Path >= ?1 and doc.Path < (?1 || '.~'))
order by (Path || '.' || Ref) asc`
const getDocumentsHierarchyQueryFilter = `
left join Document children on children.Path = (doc.Path || '.' || doc.Ref)
%[3]s
where (children.Name is null or children.Name like '%%%[1]s%%' %[2]s)
group by doc.Ref
having (doc.Name like '%%%[1]s%%' %[4]s) or count(children.Ref) > 0
order by (doc.Path || '.' || doc.Ref) asc`

type Document struct {
	Ref      string
	Path     string
	DocOrder int64
	Name     string
	NameTr   []sql.NullString
	Details  string
}

func (q *Queries) GetDocumentHierarchy(ctx context.Context, documentPath string, langs []string, filter string) ([]Document, error) {
	var query, queryTr, queryNameFilterTr, queryChildrenNameTrJoin, queryChildrenNameFilterTr strings.Builder
	var escapedFilter string
	if len(langs) > 0 {
		langs = langs[1:]
	}
	if len(filter) > 0 {
		escapedFilter = EscapeString(filter)
		query.WriteString("select doc.* from (")
	}
	query.WriteString(getDocumentsHierarchyQueryPrelude)
	args := make([]any, len(langs)+1)
	args[0] = documentPath
	for i, lang := range langs {
		args[i+1] = lang
		if len(filter) > 0 {
			fmt.Fprintf(&queryNameFilterTr, " or name_tr%d like '%%%s%%'", i+1, escapedFilter)
			fmt.Fprintf(&queryChildrenNameFilterTr, " or ctr%[1]d.Name like '%%%[2]s%%'", i+1, escapedFilter)
		}
		fmt.Fprintf(&query, ", tr%[1]d.name name_tr%[1]d", i+1)
		fmt.Fprintf(&queryTr, getDocumentsHierarchyTrQuery, i+1, i+2)
		if len(filter) > 0 {
			fmt.Fprintf(&queryChildrenNameTrJoin, getDocumentsChildrenTrQuery, i+1, i+2)
		}
	}
	query.WriteString(getDocumentsHierarchyQueryIntro)
	query.WriteString(queryTr.String())
	query.WriteString(getDocumentsHierarchyQueryCoda)
	if len(filter) > 0 {
		query.WriteString(") doc ")
		fmt.Fprintf(&query, getDocumentsHierarchyQueryFilter, escapedFilter, queryChildrenNameFilterTr.String(), queryChildrenNameTrJoin.String(), queryNameFilterTr.String())
	}
	query.WriteRune(';')
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
