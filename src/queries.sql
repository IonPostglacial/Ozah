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

-- name: GetDocumentAttachments :many
select * from Document_Attachment att 
where (att.Document_Ref = ?);

-- name: GetDocumentDirectChildren :many
select * from Document doc 
where (doc.Path = ?)
order by Path asc, Doc_Order asc;

-- name: GetDocumentHierarchy :many
select * from Document doc 
where (doc.Path like ? || '%')
order by Path asc, Doc_Order asc;
