-- name: InsertStdLangs :execresult
insert into Lang (Ref, Name) values 
    ('V', 'Vernacular'),
    ('EN', 'English'),
    ('CN', '中文'),
    ('FR', 'French'),
    ('V2', 'Vernacular Name 2'),
    ('S2', 'Name 2');

-- name: ListLangs :many
select * from Lang;

-- name: InsertDocument :execresult
insert into Document (Ref, Path, Doc_Order, Name, Details)
    values (?, ?, ?, ?, ?);

-- name: GetDocument :one
select * from Document doc where (doc.Ref = ?);

-- name: GetDocumentsNames :many
select Ref, Name from Document doc 
where doc.Ref in (sqlc.slice(path))
order by doc.Path;

-- name: GetCatCharactersNameTr2 :many
select doc.Ref, doc.Path, doc.Name, tr1.name name_tr1, tr2.name name_tr2, ch.Color from Document doc 
left join Document_Translation tr1 on doc.Ref = tr1.Document_Ref and tr1.Lang_Ref = sqlc.arg(lang1)
left join Document_Translation tr2 on doc.Ref = tr2.Document_Ref and tr2.Lang_Ref = sqlc.arg(lang2)
left join Categorical_Character ch on doc.Ref = ch.Document_Ref
where doc.Ref in (sqlc.slice(refs))
order by doc.Path, Doc_Order asc;

-- name: GetDocumentAttachments :many
select * from Document_Attachment att 
where (att.Document_Ref = ?);

-- name: GetDocumentDirectChildren :many
select * from Document doc 
where (doc.Path = ?)
order by Path asc, Doc_Order asc;

-- name: GetDocumentHierarchy :many
select * from Document doc
where (doc.Path >= sqlc.arg(path) and doc.Path < (sqlc.arg(path) || '.~'))
order by (Path || '.' || Ref) asc;

-- name: GetDocumentHierarchyTr1 :many
select doc.*, tr1.name name_tr1 from Document doc
left join Document_Translation tr1 on doc.Ref = tr1.Document_Ref and tr1.Lang_Ref = sqlc.arg(lang1)
where (doc.Path >= sqlc.arg(path) and doc.Path < (sqlc.arg(path) || '.~'))
order by (Path || '.' || Ref) asc;

-- name: GetDocumentHierarchyTr2 :many
select doc.*, tr1.name name_tr1, tr2.name name_tr2 from Document doc
left join Document_Translation tr1 on doc.Ref = tr1.Document_Ref and tr1.Lang_Ref = sqlc.arg(lang1)
left join Document_Translation tr2 on doc.Ref = tr2.Document_Ref and tr2.Lang_Ref = sqlc.arg(lang2)
where (doc.Path >= sqlc.arg(path) and doc.Path < (sqlc.arg(path) || '.~'))
order by (Path || '.' || Ref) asc;

-- name: GetTaxonInfo :one
select t.*, doc.*, tr1.name name_v, tr2.name name_cn from Taxon t
inner join Document doc on doc.Ref = t.Document_Ref
left join Document_Translation tr1 on doc.Ref = tr1.Document_Ref and tr1.Lang_Ref = "V"
left join Document_Translation tr2 on doc.Ref = tr2.Document_Ref and tr2.Lang_Ref = "CN"
where t.Document_Ref = sqlc.arg(ref);

-- name: GetDescriptors :many
select doc.Ref, doc.Path, doc.Name, doc.Details, max(att.Source) Source, count(descriptor.Description_Ref) Descriptors_Count, tr1.name name_tr1, tr2.name name_tr2 from Document doc 
left join Document_Translation tr1 on doc.Ref = tr1.Document_Ref and tr1.Lang_Ref = "EN"
left join Document_Translation tr2 on doc.Ref = tr2.Document_Ref and tr2.Lang_Ref = "CN"
left join Document_Attachment att on doc.Ref = att.Document_Ref
left join Taxon_Description descriptor on doc.Ref = descriptor.Description_Ref
where (doc.Path = ? and (descriptor.Taxon_Ref is null or descriptor.Taxon_Ref = ?))
group by doc.Ref
order by doc.Path asc, Doc_Order asc;

-- name: GetSummaryDescriptors :many
select doc.Ref, doc.Path, doc.Name, tr1.name name_tr1, tr2.name name_tr2, s.Color from Taxon_Description descriptor
inner join Document doc on doc.Ref = descriptor.Description_Ref
inner join State s on s.Document_Ref = descriptor.Description_Ref
left join Document_Translation tr1 on doc.Ref = tr1.Document_Ref and tr1.Lang_Ref = "EN"
left join Document_Translation tr2 on doc.Ref = tr2.Document_Ref and tr2.Lang_Ref = "CN"
where descriptor.Taxon_Ref = ?
order by doc.Path asc, doc.Doc_Order asc;