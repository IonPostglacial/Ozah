package documents

import (
	"fmt"
	"html/template"
	"net/http"
	"slices"
	"strconv"

	"nicolas.galipot.net/hazo/server/common"
	"nicolas.galipot.net/hazo/server/components"
	"nicolas.galipot.net/hazo/server/components/treemenu"
	"nicolas.galipot.net/hazo/server/link"
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
