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

-- name: GetDocumentDirectChildrenRefs :many
select Ref from Document doc 
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

-- name: GetTaxonStateDescriptors :many
select doc.Ref, doc.Path from Taxon_Description descriptor
inner join Document doc on doc.Ref = descriptor.Description_Ref
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

-- name: GetMeasurementCharactersWithTranslations :many
select doc.Ref, doc.Path, doc.Name, doc.Details, mc.Unit_Ref, tr1.name name_tr1, tr2.name name_tr2 from Measurement_Character mc
inner join Document doc on doc.Ref = mc.Document_Ref
left join Document_Translation tr1 on doc.Ref = tr1.Document_Ref and tr1.Lang_Ref = "EN"
left join Document_Translation tr2 on doc.Ref = tr2.Document_Ref and tr2.Lang_Ref = "CN"
order by doc.Name;

-- name: GetCategoricalCharacters :many
select doc.Ref, doc.Path, doc.Details, doc.Name, tr1.name name_tr1, tr2.name name_tr2, ch.Color from Categorical_Character ch
inner join Document doc on doc.Ref = ch.Document_Ref
left join Document_Translation tr1 on doc.Ref = tr1.Document_Ref and tr1.Lang_Ref = "EN"
left join Document_Translation tr2 on doc.Ref = tr2.Document_Ref and tr2.Lang_Ref = "CN"
order by doc.Name;

-- name: GetCharacterStates :many
select s.Document_Ref from State s
inner join Document doc on doc.Ref = s.Document_Ref
inner join Document ch on (ch.Path || '.' || ch.Ref) = doc.Path
where ch.Ref = ?;

-- name: GetTaxonBookInfo :many
select doc.Ref, doc.Name, info.Fasc, info.Page, info.Details from Taxon_Book_Info info
inner join Document doc on doc.Ref = info.Book_Ref
where info.Taxon_Ref = ?
order by doc.Name; 

-- name: InsertLang :execresult
insert into
    Lang (Ref, Name)
values
    (?, ?);

-- name: InsertDocumentTranslation :execresult
insert into Document_Translation (
    Document_Ref, 
    Lang_Ref, 
    Name, 
    Details)
values
    (?, ?, ?, ?);


-- name: InsertDocumentAttachment :execresult
insert into Document_Attachment (
    Document_Ref, 
    Attachment_Index, 
    Source, 
    Path,
    Path_Small,
    Path_Medium,
    Path_Big)
values
    (?, ?, ?, ?, ?, ?, ?);

-- name: GetDocumentAttachmentByIndex :one
select * from Document_Attachment 
where Document_Ref = ? and Attachment_Index = ?;

-- name: DeleteDocumentAttachment :exec
delete from Document_Attachment 
where Document_Ref = ? and Attachment_Index = ?;


-- name: InsertUnit :execresult
insert into
    Unit (Ref, Base_Unit_Ref, To_Base_Unit_Factor)
values
    (?, ?, ?);

-- name: InsertBook :execresult
insert into
    Book (Document_Ref, ISBN)
values
    (?, ?);

-- name: GetBooks :many
select doc.Ref, doc.Path, doc.Name, tr1.name name_tr1, tr2.name name_tr2, b.ISBN from Document doc
left join Document_Translation tr1 on doc.Ref = tr1.Document_Ref and tr1.Lang_Ref = "EN"
left join Document_Translation tr2 on doc.Ref = tr2.Document_Ref and tr2.Lang_Ref = "CN"
inner join Book b on doc.Ref = b.Document_Ref
order by doc.Path asc, doc.Doc_Order asc;

-- name: InsertState :execresult
insert into
    State (Document_Ref, Color)
values
    (?, ?);

-- name: GetAllStates :many
select doc.Ref, doc.Path, doc.Name, doc.Details, tr1.name name_tr1, tr2.name name_tr2, s.Color from Document doc
left join Document_Translation tr1 on doc.Ref = tr1.Document_Ref and tr1.Lang_Ref = "EN"
left join Document_Translation tr2 on doc.Ref = tr2.Document_Ref and tr2.Lang_Ref = "CN"
inner join State s on doc.Ref = s.Document_Ref
order by doc.Path asc, doc.Doc_Order asc;

-- name: InsertCategoricalCharacter :execresult
insert into
    Categorical_Character (Document_Ref, Color)
values
    (?, ?);

-- name: InsertMeasurementCharacter :execresult
insert into
    Measurement_Character (Document_Ref, Color, Unit_Ref)
values
    (?, ?, ?);

-- name: InsertPeriodicCharacter :execresult
insert into Periodic_Character (
    Document_Ref, 
    Periodic_Category_Ref, 
    Color)
values 
    (?, ?, ?);

-- name: InsertGeographicalPlace :execresult
insert into Geographical_Place (
    Document_Ref,
    Latitude,
    Longitude,
    Scale)
values
    (?, ?, ?, ?);

-- name: InsertGeographicalMap :execresult
insert into Geographical_Map (
    Document_Ref,
    Place_Ref,
    Map_File,
    Map_File_Feature_Name) 
values 
    (?, ?, ?, ?);

-- name: InsertGeographicalCharacter :execresult
insert into Geographical_Character (
    Document_Ref,
    Map_Ref,
	Color)
values
    (?, ?, ?);

-- name: InsertDescriptorVisibilityRequirement :execresult
insert into Descriptor_Visibility_Requirement (
    Descriptor_Ref, 
    Required_Descriptor_Ref)
values
    (?, ?);

-- name: InsertDescriptorVisibilityInapplicable :execresult
insert into Descriptor_Visibility_Inapplicable (
    Descriptor_Ref, 
    Inapplicable_Descriptor_Ref)
values
    (?, ?);

-- name: InsertTaxon :execresult
insert into Taxon (
    Document_Ref,
    Author,
    Website,
	Meaning,
    Herbarium_No,
    Herbarium_Picture,
    Fasc,
    Page)
values (?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetAllTaxonsWithTranslations :many
select t.*, doc.*, tr1.name name_v, tr2.name name_cn from Taxon t
inner join Document doc on doc.Ref = t.Document_Ref
left join Document_Translation tr1 on doc.Ref = tr1.Document_Ref and tr1.Lang_Ref = "V"
left join Document_Translation tr2 on doc.Ref = tr2.Document_Ref and tr2.Lang_Ref = "CN"
order by doc.Path asc, doc.Doc_Order asc;

-- name: InsertTaxonMeasurement :execresult
insert into Taxon_Measurement (
    Taxon_Ref ,
	Character_Ref,
	Minimum,
	Maximum) 
values (?, ?, ?, ?);

-- name: InsertTaxonDescription :execresult
insert into Taxon_Description (
    Taxon_Ref,
	Description_Ref)
values (?, ?);

-- name: InsertTaxonBookInfo :execresult
insert into Taxon_Book_Info (
    Taxon_Ref,
	Book_Ref,
	Fasc,
	Page,
    Details) 
values (?, ?, ?, ?, ?);

-- name: InsertTaxonSpecimenLocation :execresult
insert into Taxon_Specimen_Location (
    Taxon_Ref,
    Specimen_Index,
    Latitude,
    Longitude) 
values (?, ?, ?, ?);