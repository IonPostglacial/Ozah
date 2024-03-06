package server

import (
	"bytes"
	"fmt"
	"html"
	"net/http"
	"time"

	"embed"
)

//go:embed assets
var assets embed.FS

//go:embed index.html
var indexPage []byte

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
	http.ServeContent(w, r, "index.html", time.Now(), bytes.NewReader(indexPage))
}
