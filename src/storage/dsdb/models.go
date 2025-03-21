// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package dsdb

import (
	"database/sql"
)

type Book struct {
	DocumentRef string
	Isbn        sql.NullString
}

type CategoricalCharacter struct {
	DocumentRef string
	Color       sql.NullString
}

type DescriptorVisibilityInapplicable struct {
	DescriptorRef             string
	InapplicableDescriptorRef string
}

type DescriptorVisibilityRequirement struct {
	DescriptorRef         string
	RequiredDescriptorRef string
}

type Document struct {
	Ref      string
	Path     string
	DocOrder int64
	Name     string
	Details  sql.NullString
}

type DocumentAttachment struct {
	DocumentRef     string
	AttachmentIndex int64
	Source          string
	Path            string
}

type DocumentTranslation struct {
	DocumentRef string
	LangRef     string
	Name        string
	Details     sql.NullString
}

type GeographicalCharacter struct {
	DocumentRef string
	MapRef      string
	Color       sql.NullString
}

type GeographicalMap struct {
	DocumentRef        string
	PlaceRef           string
	MapFile            string
	MapFileFeatureName string
}

type GeographicalPlace struct {
	DocumentRef string
	Latitude    float64
	Longitude   float64
	Scale       int64
}

type Lang struct {
	Ref  string
	Name string
}

type MeasurementCharacter struct {
	DocumentRef string
	Color       sql.NullString
	UnitRef     sql.NullString
}

type PeriodicCharacter struct {
	DocumentRef         string
	PeriodicCategoryRef string
	Color               sql.NullString
}

type State struct {
	DocumentRef string
	Color       sql.NullString
}

type Taxon struct {
	DocumentRef      string
	Author           string
	Website          sql.NullString
	Meaning          sql.NullString
	HerbariumNo      sql.NullString
	HerbariumPicture sql.NullString
	Fasc             sql.NullInt64
	Page             sql.NullInt64
}

type TaxonBookInfo struct {
	TaxonRef string
	BookRef  string
	Fasc     sql.NullInt64
	Page     sql.NullInt64
	Details  sql.NullString
}

type TaxonDescription struct {
	TaxonRef       string
	DescriptionRef string
}

type TaxonMeasurement struct {
	TaxonRef     string
	CharacterRef string
	Minimum      sql.NullFloat64
	Maximum      sql.NullFloat64
}

type TaxonSpecimenLocation struct {
	TaxonRef      string
	SpecimenIndex int64
	Latitude      float64
	Longitude     float64
}

type Unit struct {
	Ref              string
	BaseUnitRef      sql.NullString
	ToBaseUnitFactor sql.NullFloat64
}
