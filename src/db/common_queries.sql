-- name: GetCredentials :one
select Encryption, Salt, Password, Created_On, Last_Modified from Credentials
where Login = ?;

-- name: InsertCredentials :execresult
insert into Credentials (Login, Encryption, Salt, Password)
values (?, ?, ?, ?);