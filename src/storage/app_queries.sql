-- name: GetCredentials :one
select
    Encryption,
    Password,
    Created_On,
    Last_Modified
from
    Credentials
where
    Login = ?;

-- name: InsertCredentials :execresult
insert into
    Credentials (Login, Encryption, Password)
values
    (?, ?, ?);

-- name: GetSession :one
select
    Login,
    Expiry_Date
from
    Session
where
    Token = ?;

-- name: InsertSession :execresult
insert into
    Session (Token, Login, Expiry_Date)
values
    (?, ?, ?);

-- name: DeleteUserSessions :execresult
delete from Session
where
    Login = ?;

-- name: GetUserConfiguration :one
select
    *
from
    User_Configuration
where
    Login = ?;

-- name: InsertUserConfiguration :execresult
insert into
    User_Configuration (Login, Private_Directory)
values
    (?, ?);

-- name: GetUserSelectedLanguages :many
select
    Lang_Ref
from
    User_Selected_Lang
where
    User_Login = ?;

-- name: GetUserHiddenPanels :many
select
    Panel_Id
from
    User_Hidden_Panel
where
    User_Login = ?;
