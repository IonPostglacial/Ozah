package characters

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

func (h *editActions) saveCharacter(ctx context.Context, r *http.Request) error {
	if r.PostFormValue("character-save") == "" {
		return nil
	}

	docRef := r.PostFormValue("character-ref")
	name := r.PostFormValue("nameS")
	nameEN := r.PostFormValue("nameEn")
	nameCN := r.PostFormValue("nameCn")
	description := r.PostFormValue("description")
	color := r.PostFormValue("color")
	charType := r.PostFormValue("character-type") // "categorical", "measurement", "periodic", "geographical"

	if docRef == "" {
		return fmt.Errorf("character reference is required")
	}
	if name == "" {
		return fmt.Errorf("character name is required")
	}

	_, err := h.queries.GetDocument(ctx, docRef)
	characterExists := err == nil

	if characterExists {
		err = h.queries.UpdateDocument(ctx, dsdb.UpdateDocumentParams{
			Name:     name,
			Details:  sql.NullString{String: description, Valid: true},
			DocOrder: 0,
			Ref:      docRef,
		})
		if err != nil {
			return fmt.Errorf("failed to update document: %w", err)
		}

		switch charType {
		case "categorical":
			err = h.queries.UpdateCategoricalCharacter(ctx, dsdb.UpdateCategoricalCharacterParams{
				Color:       sql.NullString{String: color, Valid: color != ""},
				DocumentRef: docRef,
			})
		case "measurement":
			unitRef := r.PostFormValue("unit-ref")
			err = h.queries.UpdateMeasurementCharacter(ctx, dsdb.UpdateMeasurementCharacterParams{
				Color:       sql.NullString{String: color, Valid: color != ""},
				UnitRef:     sql.NullString{String: unitRef, Valid: unitRef != ""},
				DocumentRef: docRef,
			})
		case "periodic":
			categoryRef := r.PostFormValue("category-ref")
			err = h.queries.UpdatePeriodicCharacter(ctx, dsdb.UpdatePeriodicCharacterParams{
				PeriodicCategoryRef: categoryRef,
				Color:               sql.NullString{String: color, Valid: color != ""},
				DocumentRef:         docRef,
			})
		case "geographical":
			mapRef := r.PostFormValue("map-ref")
			err = h.queries.UpdateGeographicalCharacter(ctx, dsdb.UpdateGeographicalCharacterParams{
				MapRef:      mapRef,
				Color:       sql.NullString{String: color, Valid: color != ""},
				DocumentRef: docRef,
			})
		}

		if err != nil {
			return fmt.Errorf("failed to update character: %w", err)
		}

		if nameEN != "" {
			err = h.queries.UpdateDocumentTranslation(ctx, dsdb.UpdateDocumentTranslationParams{
				Name:        nameEN,
				Details:     sql.NullString{},
				DocumentRef: docRef,
				LangRef:     "EN",
			})
			if err != nil {
				_, _ = h.queries.InsertDocumentTranslation(ctx, dsdb.InsertDocumentTranslationParams{
					DocumentRef: docRef,
					LangRef:     "EN",
					Name:        nameEN,
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
		return fmt.Errorf("cannot create new character: character '%s' does not exist", docRef)
	}

	return nil
}

func (h *editActions) deleteCharacter(ctx context.Context, docRef string) error {
	if docRef == "" {
		return fmt.Errorf("character reference is required")
	}

	_, err := h.queries.GetDocument(ctx, docRef)
	if err != nil {
		return fmt.Errorf("character not found: %w", err)
	}

	err = h.queries.DeleteDocument(ctx, docRef)
	if err != nil {
		return fmt.Errorf("failed to delete character: %w", err)
	}

	return nil
}

func (h *editActions) saveState(ctx context.Context, r *http.Request) error {
	if r.PostFormValue("state-save") == "" {
		return nil
	}

	docRef := r.PostFormValue("state-ref")
	name := r.PostFormValue("nameS")
	nameEN := r.PostFormValue("nameEn")
	nameCN := r.PostFormValue("nameCn")
	description := r.PostFormValue("description")
	color := r.PostFormValue("color")

	if docRef == "" {
		return fmt.Errorf("state reference is required")
	}
	if name == "" {
		return fmt.Errorf("state name is required")
	}

	_, err := h.queries.GetDocument(ctx, docRef)
	stateExists := err == nil

	if stateExists {
		err = h.queries.UpdateDocument(ctx, dsdb.UpdateDocumentParams{
			Name:     name,
			Details:  sql.NullString{String: description, Valid: true},
			DocOrder: 0,
			Ref:      docRef,
		})
		if err != nil {
			return fmt.Errorf("failed to update document: %w", err)
		}

		err = h.queries.UpdateState(ctx, dsdb.UpdateStateParams{
			Color:       sql.NullString{String: color, Valid: color != ""},
			DocumentRef: docRef,
		})
		if err != nil {
			return fmt.Errorf("failed to update state: %w", err)
		}

		if nameEN != "" {
			err = h.queries.UpdateDocumentTranslation(ctx, dsdb.UpdateDocumentTranslationParams{
				Name:        nameEN,
				Details:     sql.NullString{},
				DocumentRef: docRef,
				LangRef:     "EN",
			})
			if err != nil {
				_, _ = h.queries.InsertDocumentTranslation(ctx, dsdb.InsertDocumentTranslationParams{
					DocumentRef: docRef,
					LangRef:     "EN",
					Name:        nameEN,
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
		return fmt.Errorf("cannot create new state: state '%s' does not exist", docRef)
	}

	return nil
}

func (h *editActions) deleteState(ctx context.Context, docRef string) error {
	if docRef == "" {
		return fmt.Errorf("state reference is required")
	}

	_, err := h.queries.GetDocument(ctx, docRef)
	if err != nil {
		return fmt.Errorf("state not found: %w", err)
	}

	err = h.queries.DeleteDocument(ctx, docRef)
	if err != nil {
		return fmt.Errorf("failed to delete state: %w", err)
	}

	return nil
}

func (h *editActions) Register(reg *action.Registry) {
	reg.AppendAction(action.Action(h.saveCharacter))
	reg.AppendAction(action.NewActionWithStringArgument("character-delete", h.deleteCharacter))
	reg.AppendAction(action.Action(h.saveState))
	reg.AppendAction(action.NewActionWithStringArgument("state-delete", h.deleteState))
}
