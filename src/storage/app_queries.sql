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

-- name: InsertLang :execresult
insert into
    Lang (Ref, Name)
values
    (?, ?);

-- name: GetAllLangs :execresult
select * from Lang;

-- name: InsertUserSelectedLanguage :execresult
insert into
    User_Selected_Lang (User_Login, Lang_Ref)
values
    (?, ?);

-- name: DeleteUserSelectedLanguage :execresult
delete from User_Selected_Lang
where
    User_Login = ?
    and Lang_Ref = ?;

-- name: GetUserSelectedLanguages :many
select lang.* from Lang as lang
inner join User_Selected_Lang as selectedLang
on (lang.Ref = selectedLang.Lang_Ref)
where     
    selectedLang.User_Login = ?;

-- name: GetLangSelectionForUser :many
select 
    Ref, Name, not Lang_Ref is null as Selected
from Lang 
left join User_Selected_Lang on Ref = Lang_Ref
where 
    Lang_Ref is null or User_Login = ?;

-- name: InsertUserPanel :execresult
insert into
    Panel (Id, Name)
values
    (?, ?);

-- name: InsertUserHiddenPanels :execresult
insert into
    User_Hidden_Panel (User_Login, Panel_Id)
values
    (?, ?);

-- name: DeleteUserHiddenPanels :execresult
delete from User_Hidden_Panel
where
    User_Login = ? and Panel_Id = ?;

-- name: GetUserHiddenPanels :many
select
    Panel_Id
from
    User_Hidden_Panel
where
    User_Login = ?;
