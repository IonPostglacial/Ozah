package components

import (
	"embed"
	"html/template"
	"log"
	"slices"
	"strings"
	txtemplate "text/template"

	"nicolas.galipot.net/hazo/server/components/treemenu"
)

//go:embed */*.html
var htmlTemplates embed.FS

//go:embed */*.js
var jsFS embed.FS

var jsTemplate *txtemplate.Template
var JavascriptCode string

//go:embed */*.css
var cssFS embed.FS

var cssTemplate *txtemplate.Template
var CssCode string

func init() {
	var err error
	var buf strings.Builder
	jsTemplate, err = txtemplate.ParseFS(jsFS, "*/*.js")
	if err != nil {
		log.Fatal(err)
	}
	for _, t := range jsTemplate.Templates() {
		buf.WriteString("(()=>{")
		t.Execute(&buf, nil)
		buf.WriteString("})();")
	}
	JavascriptCode = buf.String()

	var cssbuf strings.Builder
	cssTemplate, err = txtemplate.ParseFS(cssFS, "*/*.css")
	if err != nil {
		log.Fatal(err)
	}
	for _, t := range cssTemplate.Templates() {
		t.Execute(&cssbuf, nil)
	}
	CssCode = cssbuf.String()
}

func NewTemplate() *template.Template {
	tmpl := template.New("")
	tmpl.Funcs(template.FuncMap{
		"selectedDoc": func() string {
			return ""
		},
		"sortDocs": func(items []*treemenu.Item) []*treemenu.Item {
			slices.SortFunc(items, func(i, o *treemenu.Item) int {
				return int(i.Order - o.Order)
			})
			return items
		},
		"documentUrl": func(taxon *treemenu.Item) string {
			return "#1"
		},
		"colorize": func(color string) template.HTMLAttr {
			return "style='background-color: red;'"
		},
	})
	tmpl, err := tmpl.ParseFS(htmlTemplates, "*/*.html")
	if err != nil {
		log.Fatal(err)
	}
	return tmpl
}
