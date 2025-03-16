package documents

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"slices"
	"strconv"

	"nicolas.galipot.net/hazo/server/common"
	"nicolas.galipot.net/hazo/server/components"
	"nicolas.galipot.net/hazo/server/components/treemenu"
	"nicolas.galipot.net/hazo/server/link"
	"nicolas.galipot.net/hazo/storage/appdb"
)

func HandlerWrapper(docType string) func(handler common.Handler) common.Handler {
	return func(handler common.Handler) common.Handler {
		return func(w http.ResponseWriter, r *http.Request, cc *common.Context) error {
			dsName := r.PathValue("dsName")
			docId := r.PathValue("id")
			tmpl := components.NewTemplate()
			tmpl.Funcs(template.FuncMap{
				"selectedDoc": func() string {
					return docId
				},
				"sortDocs": func(items []*treemenu.Item) []*treemenu.Item {
					slices.SortFunc(items, func(i, o *treemenu.Item) int {
						return int(i.Order - o.Order)
					})
					return items
				},
				"documentUrl": func(taxon *treemenu.Item) string {
					return link.ToDocument(dsName, taxon.Id)
				},
				"colorize": func(color string) template.HTMLAttr {
					if color == "" {
						return ""
					}
					if len(color) != 7 || color[0] != '#' {
						return ""
					}
					_, err := strconv.ParseUint(color[1:], 16, 64)
					if err != nil {
						return ""
					}
					return template.HTMLAttr(fmt.Sprintf("style='background-color: color-mix(in hsl, %s 40%%, white);'", color))
				},
			})
			cc.Template = tmpl
			err := handler(w, r, cc)
			if err == nil {
				w.Header().Add("Content-Type", "text/html")
			}
			return err
		}
	}
}

func LoadMenuLanguages(ctx context.Context, cc *common.Context, queries *appdb.Queries) ([]treemenu.Lang, []string, error) {
	langSelection, err := queries.GetLangSelectionForUser(ctx, cc.User.Login)
	if err != nil {
		return nil, nil, fmt.Errorf("couldn't retrieve the list of languages: %w", err)
	}
	menuLangs := make([]treemenu.Lang, len(langSelection)+1)
	menuSelectedLangRefs := make([]string, 1, len(langSelection)+1)
	menuLangs[0] = treemenu.Lang{Name: "S", Selected: true}
	menuSelectedLangRefs[0] = "S"
	for i, lang := range langSelection {
		menuLangs[i+1] = treemenu.Lang{
			Ref:      lang.Ref,
			Name:     lang.Name,
			Selected: lang.Selected,
		}
		if lang.Selected {
			menuSelectedLangRefs = append(menuSelectedLangRefs, lang.Ref)
		}
	}
	return menuLangs, menuSelectedLangRefs, nil
}
