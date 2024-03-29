package documents

import (
	"fmt"
	"html/template"
	"net/http"
	"slices"

	"nicolas.galipot.net/hazo/server/common"
	"nicolas.galipot.net/hazo/server/components"
	"nicolas.galipot.net/hazo/server/components/treemenu"
)

func HandlerWrapper(docType string) func(handler common.Handler) common.Handler {
	return func(handler common.Handler) common.Handler {
		return func(w http.ResponseWriter, r *http.Request, cc *common.Context) error {
			dbName := r.PathValue("dsName")
			docId := r.PathValue("id")
			tmpl := components.NewTemplate()
			tmpl.Funcs(template.FuncMap{
				"selectedDoc": func() string {
					return docId
				},
				"sortDocs": func(items []*treemenu.Item) []*treemenu.Item {
					slices.SortFunc(items, func(i, o *treemenu.Item) int {
						return i.Order - o.Order
					})
					return items
				},
				"documentUrl": func(taxon *treemenu.Item) string {
					return fmt.Sprintf("/ds/%s/%s/%s", dbName, docType, taxon.Id)
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
