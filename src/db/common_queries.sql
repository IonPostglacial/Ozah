-- name: GetCredentials :one
select Encryption, Password, Created_On, Last_Modified from Credentials
where Login = ?;

-- name: InsertCredentials :execresult
insert into Credentials (Login, Encryption, Password)
values (?, ?, ?);

-- name: GetSession :one
select Login, Expiry_Date from Session
where Token = ?;

-- name: InsertSession :execresult
insert into Session (Token, Login, Expiry_Date)
values (?, ?, ?);

-- name: DeleteUserSessions :execresult
delete from Session where Login = ?;