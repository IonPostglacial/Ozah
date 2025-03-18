package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"nicolas.galipot.net/hazo/fileformats/hazojson"
	"nicolas.galipot.net/hazo/storage/dsdb"
)

const (
	CALENDAR_ID   = "_cal"
	ROOT_PLACE_ID = "_geo"
	MADA_PLACE_ID = "_geo_mada"
)

func nullableString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: true, String: ""}
	}
	return sql.NullString{Valid: true, String: s}
}

func parseI64(s string) sql.NullInt64 {
	n, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
	if err != nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Valid: true, Int64: n}
}

func insertDocument(ctx context.Context, queries *dsdb.Queries, order int64, doc hazojson.EncodedDocument) error {
	_, err := queries.InsertDocument(ctx, dsdb.InsertDocumentParams{
		Ref: doc.Id(), Path: doc.Path(), DocOrder: order, Name: doc.Name(), Details: nullableString(doc.Description()),
	})
	if err != nil {
		return err
	}
	if len(doc.NameCN()) > 0 {
		_, err = queries.InsertDocumentTranslation(ctx, dsdb.InsertDocumentTranslationParams{
			DocumentRef: doc.Id(), LangRef: "CN", Name: doc.NameCN(), Details: sql.NullString{},
		})
		if err != nil {
			return err
		}
	}
	if len(doc.NameEN()) > 0 {
		_, err = queries.InsertDocumentTranslation(ctx, dsdb.InsertDocumentTranslationParams{
			DocumentRef: doc.Id(), LangRef: "EN", Name: doc.NameEN(), Details: sql.NullString{},
		})
		if err != nil {
			return err
		}
	}
	for i, photo := range doc.Photos() {
		_, err = queries.InsertDocumentAttachment(ctx, dsdb.InsertDocumentAttachmentParams{
			DocumentRef:     doc.Id(),
			AttachmentIndex: int64(i),
			Source:          photo.Url,
			Path:            photo.HubUrl,
		})
		if err != nil {
			return err
		}
	}
	return err
}

func WithTx(db *sql.DB, queries *dsdb.Queries, cb func(*dsdb.Queries) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := queries.WithTx(tx)
	cb(qtx)
	return tx.Commit()
}

func insertStandardData(ctx context.Context, q *dsdb.Queries) (map[string]string, error) {
	langs := []dsdb.InsertLangParams{
		{Ref: "V", Name: "Vernacular"},
		{Ref: "CN", Name: "Chinese"},
		{Ref: "EN", Name: "English"},
		{Ref: "FR", Name: "French"},
	}
	units := []dsdb.InsertUnitParams{
		{Ref: "kg", BaseUnitRef: sql.NullString{}, ToBaseUnitFactor: sql.NullFloat64{}},
		{Ref: "g", BaseUnitRef: sql.NullString{Valid: true, String: "kg"}, ToBaseUnitFactor: sql.NullFloat64{Valid: true, Float64: 1000}},
		{Ref: "m", BaseUnitRef: sql.NullString{}, ToBaseUnitFactor: sql.NullFloat64{}},
		{Ref: "mm", BaseUnitRef: sql.NullString{Valid: true, String: "m"}, ToBaseUnitFactor: sql.NullFloat64{Valid: true, Float64: 1000}},
		{Ref: "cm", BaseUnitRef: sql.NullString{Valid: true, String: "m"}, ToBaseUnitFactor: sql.NullFloat64{Valid: true, Float64: 100}},
		{Ref: "km", BaseUnitRef: sql.NullString{Valid: true, String: "m"}, ToBaseUnitFactor: sql.NullFloat64{Valid: true, Float64: 0.001}},
		{Ref: "nbr", BaseUnitRef: sql.NullString{}, ToBaseUnitFactor: sql.NullFloat64{}},
	}
	for _, lang := range langs {
		_, err := q.InsertLang(ctx, lang)
		if err != nil {
			return nil, err
		}
	}
	for _, unit := range units {
		_, err := q.InsertUnit(ctx, unit)
		if err != nil {
			return nil, err
		}
	}
	stdDocuments := []dsdb.InsertDocumentParams{
		{Ref: ROOT_PLACE_ID, Path: "", DocOrder: 0, Name: "Geographical Places", Details: sql.NullString{Valid: true, String: "All geographical places"}},
		{Ref: MADA_PLACE_ID, Path: ROOT_PLACE_ID, DocOrder: 0, Name: "Madagascar", Details: sql.NullString{Valid: true, String: "The island of Madagascar"}},
		{Ref: CALENDAR_ID, Path: "", DocOrder: 0, Name: "Calendar", Details: sql.NullString{}},
	}
	for _, doc := range stdDocuments {
		_, err := q.InsertDocument(ctx, doc)
		if err != nil {
			return nil, err
		}
	}
	_, err := q.InsertGeographicalPlace(ctx, dsdb.InsertGeographicalPlaceParams{
		DocumentRef: MADA_PLACE_ID, Latitude: -18.546564, Longitude: 46.518367, Scale: 2000,
	})
	if err != nil {
		return nil, err
	}
	stdMapsDocuments := []*hazojson.State{
		{Id: "_geo_mada_1", Name: "Province", NameEN: "Province", NameCN: "州"},
		{Id: "_geo_mada_2", Name: "Région", NameEN: "Region", NameCN: "地区"},
		{Id: "_geo_mada_3", Name: "Districte", NameEN: "District", NameCN: "区域"},
		{Id: "_geo_mada_4", Name: "Commune", NameEN: "City", NameCN: "城市"},
	}
	mapIdsByFilePath := make(map[string]string, len(stdMapsDocuments))
	for i, m := range stdMapsDocuments {
		mapIdsByFilePath[strings.Join(m.Path, ".")] = m.Id
		if err := insertDocument(ctx, q, int64(i), hazojson.StateAsDocument{State: m}); err != nil {
			return nil, err
		}
	}
	q.InsertGeographicalMap(ctx, dsdb.InsertGeographicalMapParams{
		DocumentRef: "_geo_mada_1", PlaceRef: "_geo_mada", MapFile: "MDG_adm1.json", MapFileFeatureName: "NAME_1",
	})
	q.InsertGeographicalMap(ctx, dsdb.InsertGeographicalMapParams{
		DocumentRef: "_geo_mada_2", PlaceRef: "_geo_mada", MapFile: "MDG_adm2.json", MapFileFeatureName: "NAME_2",
	})
	q.InsertGeographicalMap(ctx, dsdb.InsertGeographicalMapParams{
		DocumentRef: "_geo_mada_3", PlaceRef: "_geo_mada", MapFile: "MDG_adm3.json", MapFileFeatureName: "NAME_3",
	})
	q.InsertGeographicalMap(ctx, dsdb.InsertGeographicalMapParams{
		DocumentRef: "_geo_mada_4", PlaceRef: "_geo_mada", MapFile: "MDG_adm4.json", MapFileFeatureName: "NAME_4",
	})
	return mapIdsByFilePath, err
}

func insertBooks(ctx context.Context, data hazojson.Dataset, q *dsdb.Queries) error {
	for i, book := range data.Books {
		insertDocument(ctx, q, int64(i), hazojson.BookAsDocument{Book: book})
		_, err := q.InsertBook(ctx, dsdb.InsertBookParams{
			DocumentRef: book.Id, Isbn: sql.NullString{},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func insertStates(ctx context.Context, data hazojson.Dataset, q *dsdb.Queries) error {
	for i, state := range data.States {
		insertDocument(ctx, q, int64(i), hazojson.StateAsDocument{State: state})
		_, err := q.InsertState(ctx, dsdb.InsertStateParams{
			DocumentRef: state.Id, Color: nullableString(state.Color),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func insertCharacters(ctx context.Context, data hazojson.Dataset, mapIdsByFilePath map[string]string, q *dsdb.Queries) error {
	_, err := q.InsertDocument(ctx, dsdb.InsertDocumentParams{Ref: "c0"})
	if err != nil {
		return err
	}
	for i, character := range data.Characters {
		if err := insertDocument(ctx, q, int64(i), hazojson.CharacterAsDocument{Character: character}); err != nil {
			return err
		}
		switch character.Type {
		case hazojson.CharacterTypeDiscrete:
			for _, s := range character.InapplicableStatesIds {
				_, err = q.InsertDescriptorVisibilityInapplicable(ctx, dsdb.InsertDescriptorVisibilityInapplicableParams{
					DescriptorRef:             character.Id,
					InapplicableDescriptorRef: s,
				})
				if err != nil {
					return err
				}
			}
			for _, s := range character.RequiredStatesIds {
				_, err = q.InsertDescriptorVisibilityRequirement(ctx, dsdb.InsertDescriptorVisibilityRequirementParams{
					DescriptorRef:         character.Id,
					RequiredDescriptorRef: s,
				})
				if err != nil {
					return err
				}
			}
			switch character.Preset {
			case hazojson.CharacterPresetMap:
				_, err = q.InsertGeographicalCharacter(ctx, dsdb.InsertGeographicalCharacterParams{
					DocumentRef: character.Id, MapRef: mapIdsByFilePath[character.MapFile], Color: nullableString(character.Color),
				})
			case hazojson.CharacterPresetFlowering:
				_, err = q.InsertPeriodicCharacter(ctx, dsdb.InsertPeriodicCharacterParams{
					DocumentRef:         character.Id,
					Color:               nullableString(character.Color),
					PeriodicCategoryRef: CALENDAR_ID,
				})
			default:
				_, err = q.InsertCategoricalCharacter(ctx, dsdb.InsertCategoricalCharacterParams{
					DocumentRef: character.Id,
					Color:       nullableString(character.Color),
				})
			}
		case hazojson.CharacterTypeRange:
			_, err = q.InsertMeasurementCharacter(ctx, dsdb.InsertMeasurementCharacterParams{
				DocumentRef: character.Id,
				Color:       nullableString(character.Color),
				UnitRef:     nullableString(character.Unit),
			})
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func insertTaxons(ctx context.Context, data hazojson.Dataset, q *dsdb.Queries) error {
	_, err := q.InsertDocument(ctx, dsdb.InsertDocumentParams{Ref: "t0"})
	if err != nil {
		return err
	}
	for i, taxon := range data.Taxons {
		if err = insertDocument(ctx, q, int64(i), hazojson.TaxonAsDocument{Taxon: taxon}); err != nil {
			return err
		}
		_, err = q.InsertTaxon(ctx, dsdb.InsertTaxonParams{
			DocumentRef:      taxon.Id,
			Author:           taxon.Author,
			Website:          nullableString(taxon.Website),
			Meaning:          nullableString(taxon.Meaning),
			HerbariumNo:      nullableString(taxon.NoHerbier),
			HerbariumPicture: nullableString(taxon.HerbariumPicture),
			Fasc:             parseI64(taxon.Fasc),
			Page:             parseI64(taxon.Page),
		})
		if err != nil {
			return err
		}
		for _, m := range taxon.Measurements {
			_, err = q.InsertTaxonMeasurement(ctx, dsdb.InsertTaxonMeasurementParams{
				TaxonRef:     taxon.Id,
				CharacterRef: m.CharacterRef,
				Minimum:      sql.NullFloat64{Valid: true, Float64: m.Min},
				Maximum:      sql.NullFloat64{Valid: true, Float64: m.Max},
			})
			if err != nil {
				return err
			}
		}
		for _, d := range taxon.Descriptions {
			for _, s := range d.StatesIds {
				_, err = q.InsertTaxonDescription(ctx, dsdb.InsertTaxonDescriptionParams{
					TaxonRef:       taxon.Id,
					DescriptionRef: s,
				})
				if err != nil {
					return err
				}
			}
		}
		for id, info := range taxon.BookInfoByIds {
			_, err := q.InsertTaxonBookInfo(ctx, dsdb.InsertTaxonBookInfoParams{
				TaxonRef: taxon.Id,
				BookRef:  id,
				Details:  nullableString(info.Detail),
				Fasc:     parseI64(info.Fasc),
				Page:     parseI64(info.Page),
			})
			if err != nil {
				return err
			}
		}
		for i, loc := range taxon.SpecimenLocations {
			_, err = q.InsertTaxonSpecimenLocation(ctx, dsdb.InsertTaxonSpecimenLocationParams{
				TaxonRef:      taxon.Id,
				SpecimenIndex: int64(i),
				Latitude:      loc.Latitude,
				Longitude:     loc.Longitude,
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func ImportJson(rawData []byte, ds PrivateDataset) error {
	data := hazojson.Dataset{}
	json.Unmarshal(rawData, &data)
	db, err := ConnectDsDb(ds)
	if err != nil {
		return fmt.Errorf("importing file failed during db connection: %w", err)
	}
	ctx := context.Background()

	queries := dsdb.New(db)

	var mapIdsByFilePath map[string]string

	if err = WithTx(db, queries, func(qtx *dsdb.Queries) error {
		mapIdsByFilePath, err = insertStandardData(ctx, qtx)
		return err
	}); err != nil {
		return err
	}
	if err = WithTx(db, queries, func(qtx *dsdb.Queries) error {
		return insertBooks(ctx, data, qtx)
	}); err != nil {
		return err
	}
	if err = WithTx(db, queries, func(qtx *dsdb.Queries) error {
		return insertStates(ctx, data, qtx)
	}); err != nil {
		return err
	}
	if err = WithTx(db, queries, func(qtx *dsdb.Queries) error {
		return insertCharacters(ctx, data, mapIdsByFilePath, qtx)
	}); err != nil {
		return err
	}
	return WithTx(db, queries, func(qtx *dsdb.Queries) error {
		return insertTaxons(ctx, data, qtx)
	})
}
