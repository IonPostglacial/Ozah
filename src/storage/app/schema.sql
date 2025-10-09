create table Credentials (
    Login Text primary key,
    Encryption Text not null,
    Password Text not null,
    Created_On Text,
    Last_Modified Text
) strict,
without rowid;

create table Dataset_Sharing (
    Ref Text not null,
    Creator_User_Login Text not null,
    Creation_Date Text not null,
    Name Text not null,
    Details Text default '',
    primary key (Ref, Creator_User_Login),
    foreign key (Creator_User_Login) references Credentials (Login)
) strict,
without rowid;

create table Dataset_Sharing_Users (
    Dataset_Ref Text not null,
    Dataset_Creator_Login Text not null,
    User_Login Text not null,
    Mode Text not null,
    primary key (Dataset_Ref, Dataset_Creator_Login, User_Login),
    foreign key (Dataset_Ref, Dataset_Creator_Login) references Dataset_Sharing (Ref, Creator_User_Login),
    foreign key (User_Login) references Credentials (Login)
) strict,
without rowid;

create table User_Configuration (
    Login Text primary key,
    Private_Directory Text not null,
    foreign key (Login) references Credentials (Login)
) strict,
without rowid;

create table Lang (Ref Text primary key, Name Text not null) strict,
without rowid;

create table Panel (Id Integer primary key, Name Text not null) strict,
without rowid;

create table User_Selected_Lang (
    User_Login Text not null,
    Lang_Ref Text not null,
    primary key (User_Login, Lang_Ref),
    foreign key (User_Login) references Credentials (Login),
    foreign key (Lang_Ref) references Lang (Ref)
) strict,
without rowid;

create table User_Hidden_Panel (
    User_Login Text not null,
    Panel_Id Integer not null,
    primary key (User_Login, Panel_Id),
    foreign key (User_Login) references Credentials (Login),
    foreign key (Panel_Id) references Panel (Id)
) strict,
without rowid;

create table Session (
    Token Text primary key,
    Login Text not null,
    Expiry_Date Text not null,
    foreign key (Login) references Credentials (Login)
) strict,
without rowid;

create table Capability (
    Name Text primary key,
    Description Text not null
) strict,
without rowid;

create table User_Capability (
    User_Login Text not null,
    Capability_Name Text not null,
    Granted_Date Text not null,
    Granted_By Text not null,
    primary key (User_Login, Capability_Name),
    foreign key (User_Login) references Credentials (Login),
    foreign key (Capability_Name) references Capability (Name),
    foreign key (Granted_By) references Credentials (Login)
) strict,
without rowid;
