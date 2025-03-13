create table Credentials (
    Login Text primary key,
    Encryption Text not null,
    Password Text not null,
    Created_On Text,
    Last_Modified Text
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
