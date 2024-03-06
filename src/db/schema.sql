create table Document (
    Ref Text not null primary key,
    Path Text not null default '',
    Doc_Order Integer not null,
    Name Text not null default '',
    Details Text
);

create table Lang (
    Ref Text not null primary key,
    Name Text not null
);

create table Document_Translation (
    Document_Ref Text not null,
    Lang_Ref Text not null,

    Name Text not null,
    Details Text,

    primary key (Document_Ref, Lang_Ref),
    foreign key (Document_Ref) references Document(Ref),
    foreign key (Lang_Ref) references Lang(Ref)
);

create table Categorical_Character (
    Document_Ref Text not null primary key,
	Color Text,

    foreign key (Document_Ref) references Document(Ref) on delete cascade
);

create table Unit (
	Ref Text not null primary key,
	Base_Unit_Ref Text,
	To_Base_Unit_Factor Real,
	
	foreign key (Base_Unit_Ref) references Unit(Ref)
);

create table Mesurement_Character (
    Document_Ref Text not null primary key,
	Color Text,
	Unit_Ref Text,

    foreign key (Document_Ref) references Document(Ref) on delete cascade,
	foreign key (Unit_Ref) references Unit(Ref)
);

create table State (
    Document_Ref Text not null primary key,
	Color Text,

    foreign key (Document_Ref) references Document(Ref) on delete cascade
);

create table Document_Attachment (
    Document_Ref Text not null,
    Attachment_Index Integer not null,
	
    Source Text not null,
    Path Text not null default '',

    primary key (Document_Ref, Attachment_Index),
    foreign key (Document_Ref) references Document(Ref)
);

create table Book (
    Document_Ref Text not null primary key,
	ISBN Text,

    foreign key (Document_Ref) references Document(Ref) on delete cascade
);

create table Descriptor_Visibility_Requirement (
	Descriptor_Ref Text not null,
	Required_Descriptor_Ref Text not null,
	
	primary key (Descriptor_Ref, Required_Descriptor_Ref),
	foreign key (Descriptor_Ref) references Document(Ref),
	foreign key (Required_Descriptor_Ref) references Document(Ref)
);

create table Descriptor_Visibility_Inapplicable (
	Descriptor_Ref Text not null,
	Inapplicable_Descriptor_Ref Text not null,
	
	primary key (Descriptor_Ref, Inapplicable_Descriptor_Ref),
	foreign key (Descriptor_Ref) references Document(Ref),
	foreign key (Inapplicable_Descriptor_Ref) references Document(Ref)
);

create table Taxon (
    Document_Ref Text not null primary key,

    Author Text not null default '',
    Website Text,
	Meaning Text,
    Herbarium_No Integer,
    Herbarium_Picture Text,
    Fasc Integer,
    Page Integer,

    foreign key (Document_Ref) references Document(Ref) on delete cascade
);

create table Taxon_Measurement (
	Taxon_Ref Text not null,
	Character_Ref Text not null,
	
	Minimum Real,
	Maximum Real,
	
	primary key (Taxon_Ref, Character_Ref),
	foreign key (Taxon_Ref) references Taxon(Document_Ref),
	foreign key (Character_Ref) references Mesurement_Character(Document_Ref)
);

create table Taxon_Description (
	Taxon_Ref Text not null,
	Description_Ref Text not null,
	
	primary key (Taxon_Ref, Description_Ref),
	foreign key (Taxon_Ref) references Taxon(Document_Ref),
	foreign key (Description_Ref) references Document(Ref)
);

create table Taxon_Book_Info (
	Taxon_Ref Text not null,
	Book_Ref Text not null,

	Fasc Integer,
	Page Integer,
    Details Text,
	
	primary key (Taxon_Ref, Book_Ref),
	foreign key (Taxon_Ref) references Taxon(Document_Ref),
	foreign key (Book_Ref) references Book(Document_Ref)
);

create table Taxon_Specimen_Location (
    Taxon_Ref Text not null,
    Specimen_Index Integer not null,

    Latitude Real not null,
    Longitude Real not null,

    primary key (Taxon_Ref, Specimen_Index),
    foreign key (Taxon_Ref) references Taxon(Document_Ref)
);