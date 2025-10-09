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

-- name: UpdateCredentials :execresult
update Credentials
set
    Encryption = ?,
    Password = ?,
    Last_Modified = ?;

-- name: DeleteCredentials :execresult
delete from Credentials
where
    Login = ?;

-- name: GetDatasetSharing :many
select
    Ref,
    Creator_User_Login,
    Creation_Date,
    Name,
    Details
from
    Dataset_Sharing
where
    Creator_User_Login = ?;

-- name: InsertDatasetSharing :execresult
insert into
    Dataset_Sharing (Ref, Creator_User_Login, Creation_Date, Name, Details)
values
    (?, ?, ?, ?, ?);

-- name: UpdateDatasetSharing :execresult
update Dataset_Sharing
set
    Name = ?,
    Details = ?
where
    Ref = ?
    and Creator_User_Login = ?;

-- name: DeleteDatasetSharing :execresult
delete from Dataset_Sharing
where
    Ref = ?
    and Creator_User_Login = ?;

-- name: GetDatasetSharingUsers :many
select
    User_Login,
    Mode
from
    Dataset_Sharing_Users
where
    Dataset_Ref = ?
    and Dataset_Creator_Login = ?;

-- name: InsertDatasetSharingUser :execresult
insert into
    Dataset_Sharing_Users (Dataset_Ref, Dataset_Creator_Login, User_Login, Mode)
values
    (?, ?, ?, ?);

-- name: UpdateDatasetSharingUser :execresult
update Dataset_Sharing_Users
set
    Mode = ?
where
    Dataset_Ref = ?
    and Dataset_Creator_Login = ?
    and User_Login = ?;

-- name: DeleteDatasetSharingUser :execresult
delete from Dataset_Sharing_Users
where
    Dataset_Ref = ?
    and User_Login = ?;

-- name: GetDatasetSharingUser :one
select
    User_Login,
    Mode
from
    Dataset_Sharing_Users
where
    Dataset_Ref = ?
    and User_Login = ?;

-- name: GetReadableDatasetSharedWithUser :many
select
    ds.Ref,
    ds.Creator_User_Login,
    ds.Creation_Date,
    ds.Name,
    ds.Details,
    uc.Private_Directory
from
    Dataset_Sharing as ds
inner join
    Dataset_Sharing_Users as dsu on ds.Ref = dsu.Dataset_Ref
inner join
    User_Configuration as uc on ds.Creator_User_Login = uc.Login    
where
    dsu.User_Login = ?
    and dsu.Mode = 'read';

-- name: GetWritableDatasetSharedWithUser :many
select
    ds.Ref,
    ds.Creator_User_Login,
    ds.Creation_Date,
    ds.Name,
    ds.Details,
    uc.Private_Directory
from
    Dataset_Sharing as ds
inner join
    Dataset_Sharing_Users as dsu on ds.Ref = dsu.Dataset_Ref
inner join
    User_Configuration as uc on ds.Creator_User_Login = uc.Login
where
    dsu.User_Login = ?
    and dsu.Mode = 'write';

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

-- name: InsertCapability :execresult
insert into
    Capability (Name, Description)
values
    (?, ?);

-- name: GetUserCapabilities :many
select
    uc.Capability_Name,
    c.Description,
    uc.Granted_Date,
    uc.Granted_By
from
    User_Capability as uc
inner join
    Capability as c on uc.Capability_Name = c.Name
where
    uc.User_Login = ?;

-- name: GrantUserCapability :execresult
insert into
    User_Capability (User_Login, Capability_Name, Granted_Date, Granted_By)
values
    (?, ?, ?, ?);

-- name: RevokeUserCapability :execresult
delete from User_Capability
where
    User_Login = ?
    and Capability_Name = ?;

-- name: ListUsersWithCapability :many
select
    uc.User_Login,
    uc.Granted_Date,
    uc.Granted_By
from
    User_Capability as uc
where
    uc.Capability_Name = ?;

-- name: HasUserCapability :one
select
    count(*) > 0 as Has_Capability
from
    User_Capability
where
    User_Login = ?
    and Capability_Name = ?;
