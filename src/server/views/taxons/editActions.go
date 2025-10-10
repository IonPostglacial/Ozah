package taxons

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"nicolas.galipot.net/hazo/server/action"
	"nicolas.galipot.net/hazo/server/common"
	"nicolas.galipot.net/hazo/storage/dataset"
	"nicolas.galipot.net/hazo/storage/dsdb"
)

type editActions struct {
	cc      *common.Context
	dsName  string
	queries *dataset.Queries
}

func NewEditActions(cc *common.Context, dsName string, queries *dataset.Queries) *editActions {
	return &editActions{
		cc:      cc,
		dsName:  dsName,
		queries: queries,
	}
}

func (h *editActions) saveTaxon(ctx context.Context, r *http.Request) error {
	if r.PostFormValue("taxon-save") == "" {
		return nil
	}

	docRef := r.PostFormValue("taxon-ref")
	name := r.PostFormValue("nameS")
	author := r.PostFormValue("author")
	nameCN := r.PostFormValue("nameCn")
	nameV := r.PostFormValue("nameV")
	description := r.PostFormValue("description")
	website := r.PostFormValue("website")

	if docRef == "" {
		return fmt.Errorf("taxon reference is required")
	}
	if name == "" {
		return fmt.Errorf("taxon name is required")
	}

	_, err := h.queries.GetDocument(ctx, docRef)
	taxonExists := err == nil

	if taxonExists {
		err = h.queries.UpdateDocument(ctx, dsdb.UpdateDocumentParams{
			Name:     name,
			Details:  sql.NullString{String: description, Valid: true},
			DocOrder: 0, // Keep existing order
			Ref:      docRef,
		})
		if err != nil {
			return fmt.Errorf("failed to update document: %w", err)
		}

		err = h.queries.UpdateTaxon(ctx, dsdb.UpdateTaxonParams{
			Author:           author,
			Website:          sql.NullString{String: website, Valid: website != ""},
			Meaning:          sql.NullString{},
			HerbariumNo:      sql.NullString{},
			HerbariumPicture: sql.NullString{},
			Fasc:             sql.NullInt64{},
			Page:             sql.NullInt64{},
			DocumentRef:      docRef,
		})
		if err != nil {
			return fmt.Errorf("failed to update taxon: %w", err)
		}

		if nameV != "" {
			err = h.queries.UpdateDocumentTranslation(ctx, dsdb.UpdateDocumentTranslationParams{
				Name:        nameV,
				Details:     sql.NullString{},
				DocumentRef: docRef,
				LangRef:     "V",
			})
			if err != nil {
				_, _ = h.queries.InsertDocumentTranslation(ctx, dsdb.InsertDocumentTranslationParams{
					DocumentRef: docRef,
					LangRef:     "V",
					Name:        nameV,
					Details:     sql.NullString{},
				})
			}
		}

		if nameCN != "" {
			err = h.queries.UpdateDocumentTranslation(ctx, dsdb.UpdateDocumentTranslationParams{
				Name:        nameCN,
				Details:     sql.NullString{},
				DocumentRef: docRef,
				LangRef:     "CN",
			})
			if err != nil {
				_, _ = h.queries.InsertDocumentTranslation(ctx, dsdb.InsertDocumentTranslationParams{
					DocumentRef: docRef,
					LangRef:     "CN",
					Name:        nameCN,
					Details:     sql.NullString{},
				})
			}
		}
	} else {
		return fmt.Errorf("cannot create new taxon: taxon '%s' does not exist", docRef)
	}

	return nil
}

func (h *editActions) deleteTaxon(ctx context.Context, docRef string) error {
	if docRef == "" {
		return fmt.Errorf("taxon reference is required")
	}

	_, err := h.queries.GetTaxonInfo(ctx, docRef)
	if err != nil {
		return fmt.Errorf("taxon not found: %w", err)
	}

	err = h.queries.DeleteDocument(ctx, docRef)
	if err != nil {
		return fmt.Errorf("failed to delete taxon: %w", err)
	}

	return nil
}

func (h *editActions) saveTaxonMeasurement(ctx context.Context, r *http.Request) error {
	if r.PostFormValue("measurement-save") == "" {
		return nil
	}

	taxonRef := r.PostFormValue("taxon-ref")
	characterRef := r.PostFormValue("character-ref")
	minValue := r.PostFormValue("min-value")
	maxValue := r.PostFormValue("max-value")

	if taxonRef == "" || characterRef == "" {
		return fmt.Errorf("taxon reference and character reference are required")
	}

	var minimum, maximum sql.NullFloat64
	if minValue != "" {
		var min float64
		_, err := fmt.Sscanf(minValue, "%f", &min)
		if err != nil {
			return fmt.Errorf("invalid minimum value: %w", err)
		}
		minimum = sql.NullFloat64{Float64: min, Valid: true}
	}

	if maxValue != "" {
		var max float64
		_, err := fmt.Sscanf(maxValue, "%f", &max)
		if err != nil {
			return fmt.Errorf("invalid maximum value: %w", err)
		}
		maximum = sql.NullFloat64{Float64: max, Valid: true}
	}

	err := h.queries.UpdateTaxonMeasurement(ctx, dsdb.UpdateTaxonMeasurementParams{
		Minimum:      minimum,
		Maximum:      maximum,
		TaxonRef:     taxonRef,
		CharacterRef: characterRef,
	})
	if err != nil {
		_, err = h.queries.InsertTaxonMeasurement(ctx, dsdb.InsertTaxonMeasurementParams{
			TaxonRef:     taxonRef,
			CharacterRef: characterRef,
			Minimum:      minimum,
			Maximum:      maximum,
		})
		if err != nil {
			return fmt.Errorf("failed to save measurement: %w", err)
		}
	}

	return nil
}

func (h *editActions) deleteTaxonMeasurement(ctx context.Context, r *http.Request) error {
	if r.PostFormValue("measurement-delete") == "" {
		return nil
	}

	taxonRef := r.PostFormValue("taxon-ref")
	characterRef := r.PostFormValue("character-ref")

	if taxonRef == "" || characterRef == "" {
		return fmt.Errorf("taxon reference and character reference are required")
	}

	err := h.queries.DeleteTaxonMeasurement(ctx, dsdb.DeleteTaxonMeasurementParams{
		TaxonRef:     taxonRef,
		CharacterRef: characterRef,
	})
	if err != nil {
		return fmt.Errorf("failed to delete measurement: %w", err)
	}

	return nil
}

func (h *editActions) addTaxonDescription(ctx context.Context, r *http.Request) error {
	if r.PostFormValue("descriptor-add") == "" {
		return nil
	}

	taxonRef := r.PostFormValue("taxon-ref")
	descriptorRef := r.PostFormValue("descriptor-ref")

	if taxonRef == "" || descriptorRef == "" {
		return fmt.Errorf("taxon reference and descriptor reference are required")
	}

	_, err := h.queries.InsertTaxonDescription(ctx, dsdb.InsertTaxonDescriptionParams{
		TaxonRef:       taxonRef,
		DescriptionRef: descriptorRef,
	})
	if err != nil {
		return fmt.Errorf("failed to add descriptor: %w", err)
	}

	return nil
}

func (h *editActions) removeTaxonDescription(ctx context.Context, r *http.Request) error {
	if r.PostFormValue("descriptor-remove") == "" {
		return nil
	}

	taxonRef := r.PostFormValue("taxon-ref")
	descriptorRef := r.PostFormValue("descriptor-ref")

	if taxonRef == "" || descriptorRef == "" {
		return fmt.Errorf("taxon reference and descriptor reference are required")
	}

	err := h.queries.DeleteTaxonDescription(ctx, dsdb.DeleteTaxonDescriptionParams{
		TaxonRef:       taxonRef,
		DescriptionRef: descriptorRef,
	})
	if err != nil {
		return fmt.Errorf("failed to remove descriptor: %w", err)
	}

	return nil
}

func (h *editActions) Register(reg *action.Registry) {
	reg.AppendAction(action.Action(h.saveTaxon))
	reg.AppendAction(action.NewActionWithStringArgument("taxon-delete", h.deleteTaxon))
	reg.AppendAction(action.Action(h.saveTaxonMeasurement))
	reg.AppendAction(action.Action(h.deleteTaxonMeasurement))
	reg.AppendAction(action.Action(h.addTaxonDescription))
	reg.AppendAction(action.Action(h.removeTaxonDescription))
}
