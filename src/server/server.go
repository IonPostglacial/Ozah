package server

import (
	"context"
	"fmt"
	"html"
	"html/template"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"

	"embed"

	"nicolas.galipot.net/hazo/db"
	"nicolas.galipot.net/hazo/db/storage"
	"nicolas.galipot.net/hazo/server/components/treemenu"
)

//go:embed assets
var assets embed.FS

//go:embed index.html
var indexPage string

func Serve(addr string) error {
	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})
	http.HandleFunc("/", indexHandler)
	http.Handle("/assets/", http.FileServer(http.FS(assets)))
	http.ListenAndServe(addr, nil)
	return nil
}

type State struct {
	MenuRoot *treemenu.Item
}

func taxonHierarchyFromDb(dbName string) (*treemenu.Item, error) {
	ctx := context.Background()
	queries, err := db.Open(fmt.Sprintf("%s.sq3", dbName))
	if err != nil {
		return nil, err
	}
	docs, err := queries.GetDocumentHierarchyTr2(ctx, storage.GetDocumentHierarchyTr2Params{
		Path: "t0", Lang1: "V", Lang2: "CN",
	})
	if err != nil {
		return nil, err
	}
	h := &treemenu.Item{Id: "t0", Name: "<TOP>", FullPath: "t0"}
	previous := h
	parent := h
	breadcrumb := []*treemenu.Item{}
	for i := 0; i < len(docs); i++ {
		doc := docs[i]
		switch {
		case doc.Path == previous.FullPath:
			parent = previous
			breadcrumb = append(breadcrumb, parent)
		case doc.Path != parent.FullPath:
			for doc.Path != parent.FullPath && len(breadcrumb) > 0 {
				breadcrumb = breadcrumb[:len(breadcrumb)-1]
				parent = breadcrumb[len(breadcrumb)-1]
			}
		}
		fullPath := fmt.Sprintf("%s.%s", doc.Path, doc.Ref)
		taxon := &treemenu.Item{
			Id:       doc.Ref,
			Url:      fmt.Sprintf("/ds/%s/taxons/%s", dbName, strings.ReplaceAll(fullPath, ".", "/")),
			FullPath: fullPath,
			Order:    int(doc.DocOrder),
			Name:     doc.Name, NameV: doc.NameTr1.String, NameCN: doc.NameTr2.String,
		}
		parent.Children = append(parent.Children, taxon)
		previous = taxon
	}
	return h, nil
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.New("index")
	tmpl.Funcs(template.FuncMap{
		"sortDocs": func(items []*treemenu.Item) []*treemenu.Item {
			slices.SortFunc(items, func(i, o *treemenu.Item) int {
				return i.Order - o.Order
			})
			return items
		},
	})
	tmpl, err := tmpl.Parse(indexPage)
	if err != nil {
		log.Fatal(err)
	}
	tmpl, err = tmpl.Parse(treemenu.Template)
	if err != nil {
		log.Fatal(err)
	}
	tmpl, err = tmpl.Parse(treemenu.EntryTemplate)
	if err != nil {
		log.Fatal(err)
	}
	var buf strings.Builder
	taxons, err := taxonHierarchyFromDb("plants")
	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(&buf, State{
		MenuRoot: taxons,
	})
	if err != nil {
		log.Fatal(err)
	}
	http.ServeContent(w, r, "index.html", time.Now(), strings.NewReader(buf.String()))
}
