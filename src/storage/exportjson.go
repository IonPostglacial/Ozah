package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"nicolas.galipot.net/hazo/fileformats/hazojson"
	"nicolas.galipot.net/hazo/storage/dsdb"
)

func dbAttachmentsToPhotos(documentRef string, dbAttachments []dsdb.DocumentAttachment) []hazojson.Photo {
	photos := make([]hazojson.Photo, len(dbAttachments))
	for i, dbAttachment := range dbAttachments {
		photo := hazojson.Photo{
			Id:     fmt.Sprintf("p-%s-%d", documentRef, dbAttachment.AttachmentIndex),
			Url:    dbAttachment.Source,
			HubUrl: dbAttachment.Path,
		}
		photos[i] = photo
	}
	return photos
}

func ExportJson(dsName string, queries *Queries, w io.Writer) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	// use queries to retrieve the list of books and states
	// and put them in the dataset
	// then retrieve the tree of characters and taxons
	// and put them in the dataset
	// finally, encode the dataset as JSON
	ctx := context.Background()
	dbBooks, err := queries.GetBooks(ctx)
	if err != nil {
		return fmt.Errorf("could not get books for dataset '%s': %w", dsName, err)
	}
	books := make([]*hazojson.Book, 0, len(dbBooks))
	for _, dbBook := range dbBooks {
		book := &hazojson.Book{
			Id:    dbBook.Ref,
			Path:  strings.Split(dbBook.Path, "."),
			Label: dbBook.Name,
		}
		books = append(books, book)
	}
	dbStates, err := queries.GetAllStates(ctx)
	if err != nil {
		return fmt.Errorf("could not get states for dataset '%s': %w", dsName, err)
	}
	states := make([]*hazojson.State, 0, len(dbBooks))
	for _, dbState := range dbStates {
		dbAttachments, err := queries.GetDocumentAttachments(ctx, dbState.Ref)
		if err != nil {
			return fmt.Errorf("could not get attachments for state '%s': %w", dbState.Ref, err)
		}
		photos := dbAttachmentsToPhotos(dbState.Ref, dbAttachments)
		state := &hazojson.State{
			Id:          dbState.Ref,
			Path:        strings.Split(dbState.Path, "."),
			Name:        dbState.Name,
			NameEN:      dbState.NameTr1.String,
			NameCN:      dbState.NameTr2.String,
			Description: dbState.Details.String,
			Color:       dbState.Color.String,
			Photos:      photos,
		}
		states = append(states, state)
	}
	dbCharacters, err := queries.GetCategoricalCharacters(ctx)
	if err != nil {
		return fmt.Errorf("could not get characters for dataset '%s': %w", dsName, err)
	}
	dbMeasurementCharacters, err := queries.GetMeasurementCharactersWithTranslations(ctx)
	if err != nil {
		return fmt.Errorf("could not get measurement characters for dataset '%s': %w", dsName, err)
	}
	characters := make([]*hazojson.Character, 0, len(dbCharacters)+len(dbMeasurementCharacters))
	for _, dbCharacter := range dbCharacters {
		dbAttachments, err := queries.GetDocumentAttachments(ctx, dbCharacter.Ref)
		if err != nil {
			return fmt.Errorf("could not get attachments for character '%s': %w", dbCharacter.Ref, err)
		}
		photos := dbAttachmentsToPhotos(dbCharacter.Ref, dbAttachments)
		states, err := queries.GetCharacterStates(ctx, dbCharacter.Ref)
		if err != nil {
			return fmt.Errorf("could not get states for character '%s': %w", dbCharacter.Ref, err)
		}
		characterFullPath := FullPath(dbCharacter.Path, dbCharacter.Ref)
		childrenRefs, err := queries.GetDocumentDirectChildrenRefs(ctx, characterFullPath)
		if err != nil {
			return fmt.Errorf("could not get children for character '%s': %w", dbCharacter.Ref, err)
		}
		if childrenRefs == nil {
			childrenRefs = make([]string, 0)
		}
		if states == nil {
			states = make([]string, 0)
		}
		character := &hazojson.Character{
			Id:                    dbCharacter.Ref,
			Path:                  strings.Split(dbCharacter.Path, "."),
			Detail:                dbCharacter.Details.String,
			Name:                  dbCharacter.Name,
			NameEN:                dbCharacter.NameTr1.String,
			NameCN:                dbCharacter.NameTr2.String,
			Color:                 dbCharacter.Color.String,
			Photos:                photos,
			States:                states,
			Children:              childrenRefs,
			InapplicableStatesIds: make([]string, 0),
			RequiredStatesIds:     make([]string, 0),
		}
		characters = append(characters, character)
	}
	for _, dbMeasurementCharacter := range dbMeasurementCharacters {
		dbAttachments, err := queries.GetDocumentAttachments(ctx, dbMeasurementCharacter.Ref)
		if err != nil {
			return fmt.Errorf("could not get attachments for measurement character '%s': %w", dbMeasurementCharacter.Ref, err)
		}
		photos := dbAttachmentsToPhotos(dbMeasurementCharacter.Ref, dbAttachments)
		characters = append(characters, &hazojson.Character{
			Id:                    dbMeasurementCharacter.Ref,
			Path:                  strings.Split(dbMeasurementCharacter.Path, "."),
			Detail:                dbMeasurementCharacter.Details.String,
			Name:                  dbMeasurementCharacter.Name,
			NameEN:                dbMeasurementCharacter.NameTr1.String,
			NameCN:                dbMeasurementCharacter.NameTr2.String,
			Photos:                photos,
			States:                make([]string, 0),
			Children:              make([]string, 0),
			InapplicableStatesIds: make([]string, 0),
			RequiredStatesIds:     make([]string, 0),
			Unit:                  dbMeasurementCharacter.UnitRef.String,
		})
	}
	dbTaxons, err := queries.GetAllTaxonsWithTranslations(ctx)
	if err != nil {
		return fmt.Errorf("could not get taxons for dataset '%s': %w", dsName, err)
	}
	taxons := make([]*hazojson.Taxon, 0, len(dbTaxons))
	for _, dbTaxon := range dbTaxons {
		dbAttachments, err := queries.GetDocumentAttachments(ctx, dbTaxon.Ref)
		if err != nil {
			return fmt.Errorf("could not get attachments for taxon '%s': %w", dbTaxon.Ref, err)
		}
		photos := dbAttachmentsToPhotos(dbTaxon.Ref, dbAttachments)
		dbTaxonDescriptors, err := queries.GetTaxonStateDescriptors(ctx, dbTaxon.Ref)
		if err != nil {
			return fmt.Errorf("could not get descriptors for taxon '%s': %w", dbTaxon.Ref, err)
		}
		stateIdsByCharacterRef := make(map[string][]string)
		for _, dbDescriptor := range dbTaxonDescriptors {
			path := strings.Split(dbDescriptor.Path, ".")
			characterRef := "c0"
			if len(path) > 0 {
				characterRef = path[len(path)-1]
			}
			stateIdsByCharacterRef[characterRef] = append(stateIdsByCharacterRef[characterRef], dbDescriptor.Ref)
		}
		descriptors := make([]hazojson.Descriptions, 0, len(stateIdsByCharacterRef))
		for characterRef, stateIds := range stateIdsByCharacterRef {
			descriptors = append(descriptors, hazojson.Descriptions{
				DescriptorId: characterRef,
				StatesIds:    stateIds,
			})
		}
		taxonFullPath := FullPath(dbTaxon.Path, dbTaxon.Ref)
		childrenRefs, err := queries.GetDocumentDirectChildrenRefs(ctx, taxonFullPath)
		if err != nil {
			return fmt.Errorf("could not get children for taxon '%s': %w", dbTaxon.Ref, err)
		}
		if childrenRefs == nil {
			childrenRefs = make([]string, 0)
		}
		taxon := &hazojson.Taxon{
			Id:                dbTaxon.Ref,
			Path:              strings.Split(dbTaxon.Path, "."),
			Name:              dbTaxon.Name,
			NameEN:            dbTaxon.NameV.String,
			NameCN:            dbTaxon.NameCn.String,
			Detail:            dbTaxon.Details.String,
			Photos:            photos,
			Descriptions:      descriptors,
			Children:          childrenRefs,
			SpecimenLocations: make([]hazojson.Location, 0),
			Measurements:      make([]hazojson.Measurement, 0),
		}
		taxons = append(taxons, taxon)
	}
	ds := &hazojson.Dataset{
		Id:         dsName,
		Books:      books,
		States:     states,
		Characters: characters,
		Taxons:     taxons,
	}
	err = encoder.Encode(ds)
	if err != nil {
		return err
	}
	return nil
}
