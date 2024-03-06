package server

import (
	"fmt"
	"html"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"embed"
)

//go:embed assets
var assets embed.FS

//go:embed index.html
var indexPage string

//go:embed components/treemenu.html
var treemenu string

func Serve(addr string) error {
	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})
	http.HandleFunc("/", indexHandler)
	http.Handle("/assets/", http.FileServer(http.FS(assets)))
	http.ListenAndServe(addr, nil)
	return nil
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("index").Parse(indexPage)
	if err != nil {
		log.Fatal(err)
	}
	tmpl, err = tmpl.Parse(treemenu)
	if err != nil {
		log.Fatal(err)
	}
	var buf strings.Builder
	err = tmpl.Execute(&buf, nil)
	if err != nil {
		log.Fatal(err)
	}
	http.ServeContent(w, r, "index.html", time.Now(), strings.NewReader(buf.String()))
}
