create table Credentials (
    Login Text primary key,
    Encryption Text not null,
    Password Text not null,
    Created_On Text,
    Last_Modified Text
);

create table User_Configuration (
    Login Text primary key,
    Private_Directory text not null,
    foreign key (Login) references Credentials(Login)
);

create table Session (
    Token Text primary key,
    Login Text not null,
    Expiry_Date Text not null,

    foreign key (Login) references Credentials(Login)
);