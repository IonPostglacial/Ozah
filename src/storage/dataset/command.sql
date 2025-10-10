-- ============================================================================
-- COMMAND.SQL - Data Modification Operations (INSERT, UPDATE, DELETE)
-- ============================================================================
-- This file contains all operations that modify data in the dataset database.
-- Separated from queries.sql for clarity and maintainability.

-- ============================================================================
-- INSERT OPERATIONS
-- ============================================================================

-- name: InsertDocument :execresult
insert into Document (Ref, Path, Doc_Order, Name, Details)
    values (?, ?, ?, ?, ?);

-- name: InsertLang :execresult
insert into Lang (Ref, Name)
    values (?, ?);

-- name: InsertDocumentTranslation :execresult
insert into Document_Translation (Document_Ref, Lang_Ref, Name, Details)
    values (?, ?, ?, ?);

-- name: InsertDocumentAttachment :execresult
insert into Document_Attachment (
    Document_Ref, 
    Attachment_Index, 
    Source, 
    Path,
    Path_Small,
    Path_Medium,
    Path_Big)
values (?, ?, ?, ?, ?, ?, ?);

-- name: InsertUnit :execresult
insert into Unit (Ref, Base_Unit_Ref, To_Base_Unit_Factor)
    values (?, ?, ?);

-- name: InsertBook :execresult
insert into Book (Document_Ref, ISBN)
    values (?, ?);

-- name: InsertState :execresult
insert into State (Document_Ref, Color)
    values (?, ?);

-- name: InsertCategoricalCharacter :execresult
insert into Categorical_Character (Document_Ref, Color)
    values (?, ?);

-- name: InsertMeasurementCharacter :execresult
insert into Measurement_Character (Document_Ref, Color, Unit_Ref)
    values (?, ?, ?);

-- name: InsertPeriodicCharacter :execresult
insert into Periodic_Character (Document_Ref, Periodic_Category_Ref, Color)
    values (?, ?, ?);

-- name: InsertGeographicalPlace :execresult
insert into Geographical_Place (Document_Ref, Latitude, Longitude, Scale)
    values (?, ?, ?, ?);

-- name: InsertGeographicalMap :execresult
insert into Geographical_Map (Document_Ref, Place_Ref, Map_File, Map_File_Feature_Name) 
    values (?, ?, ?, ?);

-- name: InsertGeographicalCharacter :execresult
insert into Geographical_Character (Document_Ref, Map_Ref, Color)
    values (?, ?, ?);

-- name: InsertDescriptorVisibilityRequirement :execresult
insert into Descriptor_Visibility_Requirement (Descriptor_Ref, Required_Descriptor_Ref)
    values (?, ?);

-- name: InsertDescriptorVisibilityInapplicable :execresult
insert into Descriptor_Visibility_Inapplicable (Descriptor_Ref, Inapplicable_Descriptor_Ref)
    values (?, ?);

-- name: InsertTaxon :execresult
insert into Taxon (
    Document_Ref,
    Author,
    Website,
    Meaning,
    Herbarium_No,
    Herbarium_Picture,
    Fasc,
    Page)
values (?, ?, ?, ?, ?, ?, ?, ?);

-- name: InsertTaxonMeasurement :execresult
insert into Taxon_Measurement (Taxon_Ref, Character_Ref, Minimum, Maximum) 
    values (?, ?, ?, ?);

-- name: InsertTaxonDescription :execresult
insert into Taxon_Description (Taxon_Ref, Description_Ref)
    values (?, ?);

-- name: InsertTaxonBookInfo :execresult
insert into Taxon_Book_Info (Taxon_Ref, Book_Ref, Fasc, Page, Details) 
    values (?, ?, ?, ?, ?);

-- name: InsertTaxonSpecimenLocation :execresult
insert into Taxon_Specimen_Location (Taxon_Ref, Specimen_Index, Latitude, Longitude) 
    values (?, ?, ?, ?);

-- ============================================================================
-- UPDATE OPERATIONS
-- ============================================================================

-- name: UpdateDocument :exec
update Document 
set Name = ?, Details = ?, Doc_Order = ?
where Ref = ?;

-- name: UpdateDocumentTranslation :exec
update Document_Translation 
set Name = ?, Details = ?
where Document_Ref = ? and Lang_Ref = ?;

-- name: UpdateTaxon :exec
update Taxon 
set Author = ?, Website = ?, Meaning = ?, Herbarium_No = ?, 
    Herbarium_Picture = ?, Fasc = ?, Page = ?
where Document_Ref = ?;

-- name: UpdateCategoricalCharacter :exec
update Categorical_Character 
set Color = ?
where Document_Ref = ?;

-- name: UpdateMeasurementCharacter :exec
update Measurement_Character 
set Color = ?, Unit_Ref = ?
where Document_Ref = ?;

-- name: UpdatePeriodicCharacter :exec
update Periodic_Character 
set Periodic_Category_Ref = ?, Color = ?
where Document_Ref = ?;

-- name: UpdateGeographicalCharacter :exec
update Geographical_Character 
set Map_Ref = ?, Color = ?
where Document_Ref = ?;

-- name: UpdateState :exec
update State 
set Color = ?
where Document_Ref = ?;

-- name: UpdateTaxonMeasurement :exec
update Taxon_Measurement 
set Minimum = ?, Maximum = ?
where Taxon_Ref = ? and Character_Ref = ?;

-- name: UpdateBook :exec
update Book 
set ISBN = ?
where Document_Ref = ?;

-- name: UpdateUnit :exec
update Unit 
set Base_Unit_Ref = ?, To_Base_Unit_Factor = ?
where Ref = ?;

-- name: UpdateGeographicalPlace :exec
update Geographical_Place 
set Latitude = ?, Longitude = ?, Scale = ?
where Document_Ref = ?;

-- name: UpdateGeographicalMap :exec
update Geographical_Map 
set Place_Ref = ?, Map_File = ?, Map_File_Feature_Name = ?
where Document_Ref = ?;

-- ============================================================================
-- DELETE OPERATIONS
-- ============================================================================

-- name: DeleteDocument :exec
delete from Document where Ref = ?;

-- name: DeleteDocumentTranslation :exec
delete from Document_Translation 
where Document_Ref = ? and Lang_Ref = ?;

-- name: DeleteDocumentAttachment :exec
delete from Document_Attachment 
where Document_Ref = ? and Attachment_Index = ?;

-- name: DeleteTaxon :exec
delete from Taxon where Document_Ref = ?;

-- name: DeleteCategoricalCharacter :exec
delete from Categorical_Character where Document_Ref = ?;

-- name: DeleteMeasurementCharacter :exec
delete from Measurement_Character where Document_Ref = ?;

-- name: DeletePeriodicCharacter :exec
delete from Periodic_Character where Document_Ref = ?;

-- name: DeleteGeographicalCharacter :exec
delete from Geographical_Character where Document_Ref = ?;

-- name: DeleteState :exec
delete from State where Document_Ref = ?;

-- name: DeleteTaxonMeasurement :exec
delete from Taxon_Measurement 
where Taxon_Ref = ? and Character_Ref = ?;

-- name: DeleteTaxonDescription :exec
delete from Taxon_Description 
where Taxon_Ref = ? and Description_Ref = ?;

-- name: DeleteTaxonBookInfo :exec
delete from Taxon_Book_Info 
where Taxon_Ref = ? and Book_Ref = ?;

-- name: DeleteTaxonSpecimenLocation :exec
delete from Taxon_Specimen_Location 
where Taxon_Ref = ? and Specimen_Index = ?;

-- name: DeleteBook :exec
delete from Book where Document_Ref = ?;

-- name: DeleteUnit :exec
delete from Unit where Ref = ?;

-- name: DeleteLang :exec
delete from Lang where Ref = ?;

-- name: DeleteGeographicalPlace :exec
delete from Geographical_Place where Document_Ref = ?;

-- name: DeleteGeographicalMap :exec
delete from Geographical_Map where Document_Ref = ?;

-- name: DeleteDescriptorVisibilityRequirement :exec
delete from Descriptor_Visibility_Requirement 
where Descriptor_Ref = ? and Required_Descriptor_Ref = ?;

-- name: DeleteDescriptorVisibilityInapplicable :exec
delete from Descriptor_Visibility_Inapplicable 
where Descriptor_Ref = ? and Inapplicable_Descriptor_Ref = ?;
