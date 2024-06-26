-- name: ListLangs :many
select lang.* from Lang lang;

-- name: InsertDocument :execresult
insert into Document (Ref, Path, Doc_Order, Name, Details)
    values (?, ?, ?, ?, ?);

-- name: GetDocument :one
select * from Document doc where (doc.Ref = ?);

-- name: GetDocumentTr2 :one
select doc.*, tr1.name name_tr1, tr2.name name_tr2 from Document doc
left join Document_Translation tr1 on doc.Ref = tr1.Document_Ref and tr1.Lang_Ref = sqlc.arg(lang1)
left join Document_Translation tr2 on doc.Ref = tr2.Document_Ref and tr2.Lang_Ref = sqlc.arg(lang2)
where (doc.Ref = sqlc.arg(ref));

-- name: GetDocumentsNames :many
select Ref, Name from Document doc 
where doc.Ref in (sqlc.slice(refs))
order by doc.Path;

-- name: GetDocumentsAndParentsNames :many
select parent.Ref Parent_Ref, parent.Name Parent_Name, doc.Ref, doc.Name from Document doc
left join Document parent on doc.Path = (parent.Path || '.' || parent.Ref)
where doc.Ref in (sqlc.slice(refs))
order by parent.Ref, doc.Doc_Order;

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

-- name: GetTaxonInfo :one
select t.*, doc.*, tr1.name name_v, tr2.name name_cn from Taxon t
inner join Document doc on doc.Ref = t.Document_Ref
left join Document_Translation tr1 on doc.Ref = tr1.Document_Ref and tr1.Lang_Ref = "V"
left join Document_Translation tr2 on doc.Ref = tr2.Document_Ref and tr2.Lang_Ref = "CN"
where t.Document_Ref = sqlc.arg(ref);

-- name: GetDescriptors :many
select 
    doc.Ref, doc.Path, doc.Name, doc.Details, 
    att.Source, (descriptor.Taxon_Ref is null) Unselected, 
    tr1.name name_tr1, tr2.name name_tr2 
from Document doc 
left join Document_Translation tr1 on doc.Ref = tr1.Document_Ref and tr1.Lang_Ref = "EN"
left join Document_Translation tr2 on doc.Ref = tr2.Document_Ref and tr2.Lang_Ref = "CN"
left join Document_Attachment att on doc.Ref = att.Document_Ref
left join Taxon_Description descriptor on doc.Ref = descriptor.Description_Ref
where (doc.Path = ? and (descriptor.Taxon_Ref is null or descriptor.Taxon_Ref = ?))
order by doc.Path asc, Doc_Order asc;

-- name: GetSummaryDescriptors :many
select doc.Ref, doc.Path, doc.Name, tr1.name name_tr1, tr2.name name_tr2, s.Color from Taxon_Description descriptor
inner join Document doc on doc.Ref = descriptor.Description_Ref
inner join State s on s.Document_Ref = descriptor.Description_Ref
left join Document_Translation tr1 on doc.Ref = tr1.Document_Ref and tr1.Lang_Ref = "EN"
left join Document_Translation tr2 on doc.Ref = tr2.Document_Ref and tr2.Lang_Ref = "CN"
where descriptor.Taxon_Ref = ?
order by doc.Path asc, doc.Doc_Order asc;

-- name: DistinctiveCharacters :many
select ch.Ref, ch.Name, s.Ref State_ref, s.Name State_Name from Document ch
inner join Document s on s.Path = (ch.Path || '.' || ch.Ref)
where (ch.Path || '.' || ch.Ref) in (
    select doc.Path from Document doc 
    where Ref in (
        select Description_Ref from Taxon_Description 
        group by Description_ref 
        order by count(Taxon_Ref)
    )
);

-- name: GetMeasurementCharacters :many
select doc.Ref, doc.Name, mc.Unit_Ref from Measurement_Character mc
inner join Document doc on doc.Ref = mc.Document_Ref
order by doc.Name;

-- name: GetTaxonBookInfo :many
select doc.Ref, doc.Name, info.Fasc, info.Page, info.Details from Taxon_Book_Info info
inner join Document doc on doc.Ref = info.Book_Ref
where info.Taxon_Ref = ?
order by doc.Name; 