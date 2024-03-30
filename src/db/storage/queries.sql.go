// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: queries.sql

package storage

import (
	"context"
	"database/sql"
	"strings"
)

const getCatCharactersNameTr2 = `-- name: GetCatCharactersNameTr2 :many
select doc.Ref, doc.Path, doc.Name, tr1.name name_tr1, tr2.name name_tr2, ch.Color from Document doc 
left join Document_Translation tr1 on doc.Ref = tr1.Document_Ref and tr1.Lang_Ref = ?
left join Document_Translation tr2 on doc.Ref = tr2.Document_Ref and tr2.Lang_Ref = ?
left join Categorical_Character ch on doc.Ref = ch.Document_Ref
where doc.Ref in (/*SLICE:refs*/?)
order by doc.Path, Doc_Order asc
`

type GetCatCharactersNameTr2Params struct {
	Lang1 string
	Lang2 string
	Refs  []string
}

type GetCatCharactersNameTr2Row struct {
	Ref     string
	Path    string
	Name    string
	NameTr1 sql.NullString
	NameTr2 sql.NullString
	Color   sql.NullString
}

func (q *Queries) GetCatCharactersNameTr2(ctx context.Context, arg GetCatCharactersNameTr2Params) ([]GetCatCharactersNameTr2Row, error) {
	query := getCatCharactersNameTr2
	var queryParams []interface{}
	queryParams = append(queryParams, arg.Lang1)
	queryParams = append(queryParams, arg.Lang2)
	if len(arg.Refs) > 0 {
		for _, v := range arg.Refs {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:refs*/?", strings.Repeat(",?", len(arg.Refs))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:refs*/?", "NULL", 1)
	}
	rows, err := q.db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCatCharactersNameTr2Row
	for rows.Next() {
		var i GetCatCharactersNameTr2Row
		if err := rows.Scan(
			&i.Ref,
			&i.Path,
			&i.Name,
			&i.NameTr1,
			&i.NameTr2,
			&i.Color,
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

const getDescriptors = `-- name: GetDescriptors :many
select doc.Ref, doc.Path, doc.Name, doc.Details, max(att.Source) Source, count(descriptor.Description_Ref) Descriptors_Count, tr1.name name_tr1, tr2.name name_tr2 from Document doc 
left join Document_Translation tr1 on doc.Ref = tr1.Document_Ref and tr1.Lang_Ref = "EN"
left join Document_Translation tr2 on doc.Ref = tr2.Document_Ref and tr2.Lang_Ref = "CN"
left join Document_Attachment att on doc.Ref = att.Document_Ref
left join Taxon_Description descriptor on doc.Ref = descriptor.Description_Ref
where (doc.Path = ? and (descriptor.Taxon_Ref is null or descriptor.Taxon_Ref = ?))
group by doc.Ref
order by doc.Path asc, Doc_Order asc
`

type GetDescriptorsParams struct {
	Path     string
	TaxonRef string
}

type GetDescriptorsRow struct {
	Ref              string
	Path             string
	Name             string
	Details          sql.NullString
	Source           interface{}
	DescriptorsCount int64
	NameTr1          sql.NullString
	NameTr2          sql.NullString
}

func (q *Queries) GetDescriptors(ctx context.Context, arg GetDescriptorsParams) ([]GetDescriptorsRow, error) {
	rows, err := q.db.QueryContext(ctx, getDescriptors, arg.Path, arg.TaxonRef)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetDescriptorsRow
	for rows.Next() {
		var i GetDescriptorsRow
		if err := rows.Scan(
			&i.Ref,
			&i.Path,
			&i.Name,
			&i.Details,
			&i.Source,
			&i.DescriptorsCount,
			&i.NameTr1,
			&i.NameTr2,
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

const getDocument = `-- name: GetDocument :one
select ref, path, doc_order, name, details from Document doc where (doc.Ref = ?)
`

func (q *Queries) GetDocument(ctx context.Context, ref string) (Document, error) {
	row := q.db.QueryRowContext(ctx, getDocument, ref)
	var i Document
	err := row.Scan(
		&i.Ref,
		&i.Path,
		&i.DocOrder,
		&i.Name,
		&i.Details,
	)
	return i, err
}

const getDocumentAttachments = `-- name: GetDocumentAttachments :many
select document_ref, attachment_index, source, path from Document_Attachment att 
where (att.Document_Ref = ?)
`

func (q *Queries) GetDocumentAttachments(ctx context.Context, documentRef string) ([]DocumentAttachment, error) {
	rows, err := q.db.QueryContext(ctx, getDocumentAttachments, documentRef)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []DocumentAttachment
	for rows.Next() {
		var i DocumentAttachment
		if err := rows.Scan(
			&i.DocumentRef,
			&i.AttachmentIndex,
			&i.Source,
			&i.Path,
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

const getDocumentDirectChildren = `-- name: GetDocumentDirectChildren :many
select ref, path, doc_order, name, details from Document doc 
where (doc.Path = ?)
order by Path asc, Doc_Order asc
`

func (q *Queries) GetDocumentDirectChildren(ctx context.Context, path string) ([]Document, error) {
	rows, err := q.db.QueryContext(ctx, getDocumentDirectChildren, path)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Document
	for rows.Next() {
		var i Document
		if err := rows.Scan(
			&i.Ref,
			&i.Path,
			&i.DocOrder,
			&i.Name,
			&i.Details,
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

const getDocumentHierarchy = `-- name: GetDocumentHierarchy :many
select ref, path, doc_order, name, details from Document doc
where (doc.Path >= ? and doc.Path < (? || '.~'))
order by (Path || '.' || Ref) asc
`

type GetDocumentHierarchyParams struct {
	Path string
}

func (q *Queries) GetDocumentHierarchy(ctx context.Context, arg GetDocumentHierarchyParams) ([]Document, error) {
	rows, err := q.db.QueryContext(ctx, getDocumentHierarchy, arg.Path, arg.Path)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Document
	for rows.Next() {
		var i Document
		if err := rows.Scan(
			&i.Ref,
			&i.Path,
			&i.DocOrder,
			&i.Name,
			&i.Details,
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

const getDocumentHierarchyTr1 = `-- name: GetDocumentHierarchyTr1 :many
select doc.ref, doc.path, doc.doc_order, doc.name, doc.details, tr1.name name_tr1 from Document doc
left join Document_Translation tr1 on doc.Ref = tr1.Document_Ref and tr1.Lang_Ref = ?
where (doc.Path >= ? and doc.Path < (? || '.~'))
order by (Path || '.' || Ref) asc
`

type GetDocumentHierarchyTr1Params struct {
	Lang1 string
	Path  string
}

type GetDocumentHierarchyTr1Row struct {
	Ref      string
	Path     string
	DocOrder int32
	Name     string
	Details  sql.NullString
	NameTr1  sql.NullString
}

func (q *Queries) GetDocumentHierarchyTr1(ctx context.Context, arg GetDocumentHierarchyTr1Params) ([]GetDocumentHierarchyTr1Row, error) {
	rows, err := q.db.QueryContext(ctx, getDocumentHierarchyTr1, arg.Lang1, arg.Path, arg.Path)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetDocumentHierarchyTr1Row
	for rows.Next() {
		var i GetDocumentHierarchyTr1Row
		if err := rows.Scan(
			&i.Ref,
			&i.Path,
			&i.DocOrder,
			&i.Name,
			&i.Details,
			&i.NameTr1,
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

const getDocumentHierarchyTr2 = `-- name: GetDocumentHierarchyTr2 :many
select doc.ref, doc.path, doc.doc_order, doc.name, doc.details, tr1.name name_tr1, tr2.name name_tr2 from Document doc
left join Document_Translation tr1 on doc.Ref = tr1.Document_Ref and tr1.Lang_Ref = ?
left join Document_Translation tr2 on doc.Ref = tr2.Document_Ref and tr2.Lang_Ref = ?
where (doc.Path >= ? and doc.Path < (? || '.~'))
order by (Path || '.' || Ref) asc
`

type GetDocumentHierarchyTr2Params struct {
	Lang1 string
	Lang2 string
	Path  string
}

type GetDocumentHierarchyTr2Row struct {
	Ref      string
	Path     string
	DocOrder int32
	Name     string
	Details  sql.NullString
	NameTr1  sql.NullString
	NameTr2  sql.NullString
}

func (q *Queries) GetDocumentHierarchyTr2(ctx context.Context, arg GetDocumentHierarchyTr2Params) ([]GetDocumentHierarchyTr2Row, error) {
	rows, err := q.db.QueryContext(ctx, getDocumentHierarchyTr2,
		arg.Lang1,
		arg.Lang2,
		arg.Path,
		arg.Path,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetDocumentHierarchyTr2Row
	for rows.Next() {
		var i GetDocumentHierarchyTr2Row
		if err := rows.Scan(
			&i.Ref,
			&i.Path,
			&i.DocOrder,
			&i.Name,
			&i.Details,
			&i.NameTr1,
			&i.NameTr2,
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

const getDocumentsNames = `-- name: GetDocumentsNames :many
select Ref, Name from Document doc 
where doc.Ref in (/*SLICE:path*/?)
order by doc.Path
`

type GetDocumentsNamesRow struct {
	Ref  string
	Name string
}

func (q *Queries) GetDocumentsNames(ctx context.Context, path []string) ([]GetDocumentsNamesRow, error) {
	query := getDocumentsNames
	var queryParams []interface{}
	if len(path) > 0 {
		for _, v := range path {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:path*/?", strings.Repeat(",?", len(path))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:path*/?", "NULL", 1)
	}
	rows, err := q.db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetDocumentsNamesRow
	for rows.Next() {
		var i GetDocumentsNamesRow
		if err := rows.Scan(&i.Ref, &i.Name); err != nil {
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

const getSummaryDescriptors = `-- name: GetSummaryDescriptors :many
select doc.Ref, doc.Path, doc.Name, tr1.name name_tr1, tr2.name name_tr2, s.Color from Taxon_Description descriptor
inner join Document doc on doc.Ref = descriptor.Description_Ref
inner join State s on s.Document_Ref = descriptor.Description_Ref
left join Document_Translation tr1 on doc.Ref = tr1.Document_Ref and tr1.Lang_Ref = "EN"
left join Document_Translation tr2 on doc.Ref = tr2.Document_Ref and tr2.Lang_Ref = "CN"
where descriptor.Taxon_Ref = ?
order by doc.Path asc, doc.Doc_Order asc
`

type GetSummaryDescriptorsRow struct {
	Ref     string
	Path    string
	Name    string
	NameTr1 sql.NullString
	NameTr2 sql.NullString
	Color   sql.NullString
}

func (q *Queries) GetSummaryDescriptors(ctx context.Context, taxonRef string) ([]GetSummaryDescriptorsRow, error) {
	rows, err := q.db.QueryContext(ctx, getSummaryDescriptors, taxonRef)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetSummaryDescriptorsRow
	for rows.Next() {
		var i GetSummaryDescriptorsRow
		if err := rows.Scan(
			&i.Ref,
			&i.Path,
			&i.Name,
			&i.NameTr1,
			&i.NameTr2,
			&i.Color,
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

const getTaxonInfo = `-- name: GetTaxonInfo :one
select t.document_ref, t.author, t.website, t.meaning, t.herbarium_no, t.herbarium_picture, t.fasc, t.page, doc.ref, doc.path, doc.doc_order, doc.name, doc.details, tr1.name name_v, tr2.name name_cn from Taxon t
inner join Document doc on doc.Ref = t.Document_Ref
left join Document_Translation tr1 on doc.Ref = tr1.Document_Ref and tr1.Lang_Ref = "V"
left join Document_Translation tr2 on doc.Ref = tr2.Document_Ref and tr2.Lang_Ref = "CN"
where t.Document_Ref = ?
`

type GetTaxonInfoRow struct {
	DocumentRef      string
	Author           string
	Website          sql.NullString
	Meaning          sql.NullString
	HerbariumNo      sql.NullString
	HerbariumPicture sql.NullString
	Fasc             sql.NullString
	Page             sql.NullString
	Ref              string
	Path             string
	DocOrder         int32
	Name             string
	Details          sql.NullString
	NameV            sql.NullString
	NameCn           sql.NullString
}

func (q *Queries) GetTaxonInfo(ctx context.Context, ref string) (GetTaxonInfoRow, error) {
	row := q.db.QueryRowContext(ctx, getTaxonInfo, ref)
	var i GetTaxonInfoRow
	err := row.Scan(
		&i.DocumentRef,
		&i.Author,
		&i.Website,
		&i.Meaning,
		&i.HerbariumNo,
		&i.HerbariumPicture,
		&i.Fasc,
		&i.Page,
		&i.Ref,
		&i.Path,
		&i.DocOrder,
		&i.Name,
		&i.Details,
		&i.NameV,
		&i.NameCn,
	)
	return i, err
}

const insertDocument = `-- name: InsertDocument :execresult
insert into Document (Ref, Path, Doc_Order, Name, Details)
    values (?, ?, ?, ?, ?)
`

type InsertDocumentParams struct {
	Ref      string
	Path     string
	DocOrder int32
	Name     string
	Details  sql.NullString
}

func (q *Queries) InsertDocument(ctx context.Context, arg InsertDocumentParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, insertDocument,
		arg.Ref,
		arg.Path,
		arg.DocOrder,
		arg.Name,
		arg.Details,
	)
}

const insertStdLangs = `-- name: InsertStdLangs :execresult
insert into Lang (Ref, Name) values 
    ('V', 'Vernacular'),
    ('EN', 'English'),
    ('CN', '中文'),
    ('FR', 'French'),
    ('V2', 'Vernacular Name 2'),
    ('S2', 'Name 2')
`

func (q *Queries) InsertStdLangs(ctx context.Context) (sql.Result, error) {
	return q.db.ExecContext(ctx, insertStdLangs)
}

const listLangs = `-- name: ListLangs :many
select ref, name from Lang
`

func (q *Queries) ListLangs(ctx context.Context) ([]Lang, error) {
	rows, err := q.db.QueryContext(ctx, listLangs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Lang
	for rows.Next() {
		var i Lang
		if err := rows.Scan(&i.Ref, &i.Name); err != nil {
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
