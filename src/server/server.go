package server

import (
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"embed"

	"nicolas.galipot.net/hazo/db"
	"nicolas.galipot.net/hazo/server/authentication"
	"nicolas.galipot.net/hazo/server/common"
	"nicolas.galipot.net/hazo/server/views/taxons"
)

//go:embed assets
var assets embed.FS

//go:embed index.html
var indexPage string

func New() *http.ServeMux {
	s := http.NewServeMux()
	s.HandleFunc("/ds/{dsName}/taxons", authentication.HandlerWrapper(taxons.Handler))
	s.HandleFunc("/ds/{dsName}/taxons/{id}", authentication.HandlerWrapper(taxons.Handler))
	s.HandleFunc("/", authentication.HandlerWrapper(indexHandler))
	s.Handle("/assets/", http.FileServer(http.FS(assets)))
	return s
}

type State struct {
	Datasets []db.Dataset
}

func indexHandler(w http.ResponseWriter, r *http.Request, cc *common.Context) {
	tmpl := template.New("index")
	tmpl, err := tmpl.Parse(indexPage)
	if err != nil {
		log.Fatal(err)
	}
	datasets, err := db.ListDatasets()
	if err != nil {
		log.Fatal(err)
	}
	var buf strings.Builder
	err = tmpl.Execute(&buf, State{
		Datasets: datasets,
	})
	if err != nil {
		log.Fatal(err)
	}
	http.ServeContent(w, r, "index.html", time.Now(), strings.NewReader(buf.String()))
}
