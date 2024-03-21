create table Credentials (
    Login Text primary key,
    Encryption Text not null,
    Password Text not null,
    Created_On Text,
    Last_Modified Text
);