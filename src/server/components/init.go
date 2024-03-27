package components

import (
	"embed"
	"html/template"
	"log"
	"slices"

	"nicolas.galipot.net/hazo/server/components/treemenu"
)

//go:embed */*.html
var htmlTemplates embed.FS

func NewTemplate() *template.Template {
	tmpl := template.New("")
	tmpl.Funcs(template.FuncMap{
		"selectedDoc": func() string {
			return ""
		},
		"sortDocs": func(items []*treemenu.Item) []*treemenu.Item {
			slices.SortFunc(items, func(i, o *treemenu.Item) int {
				return i.Order - o.Order
			})
			return items
		},
		"documentUrl": func(taxon *treemenu.Item) string {
			return "#1"
		},
	})
	tmpl, err := tmpl.ParseFS(htmlTemplates, "*/*.html")
	if err != nil {
		log.Fatal("fuck", err)
	}
	return tmpl
}
