package server

import (
	"archive/zip"
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

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
	s.HandleFunc("/upload", authentication.HandlerWrapper(uploadHandler))
	s.HandleFunc("/", authentication.HandlerWrapper(indexHandler))
	s.Handle("/assets/", http.FileServer(http.FS(assets)))
	return s
}

type State struct {
	Datasets []db.Dataset
}

func uploadHandler(w http.ResponseWriter, r *http.Request, cc *common.Context) {
	fmt.Println("handle")
	r.ParseMultipartForm(5000000) //5 MB in memory, the rest in disk
	datas := r.MultipartForm
	for _, headers := range datas.File {
		auxiliar, _ := headers[0].Open() //TODO: first check len(headers) is correct
		fileName := headers[0].Filename
		dir, err := os.MkdirTemp("tmp", fileName)
		if err != nil {
			log.Fatal(err)
		}
		defer os.RemoveAll(dir)
		file, _ := io.ReadAll(auxiliar)
		r, err := zip.NewReader(bytes.NewReader(file), int64(len(file)))
		if err != nil {
			fmt.Printf("error reading zip: %s\n", err)
		}
		for _, f := range r.File {
			content, err := f.Open()
			if err != nil {
				log.Fatalf("error reading file '%s': %s\n", f.Name, err)
			}
			filePath := path.Join(dir, f.Name)
			file, err := os.Create(filePath)
			if err != nil {
				log.Fatalf("error creating file '%s': %s\n", filePath, err)
			}
			io.Copy(file, content)
			dbPath := fmt.Sprintf("%s.sq3", strings.TrimSuffix(fileName, ".zip"))
			err = db.Init(dbPath)
			if err != nil {
				log.Fatalf("error creating db '%s': %s\n", dbPath, err)
			}
			err = db.ImportCsv(dir, dbPath)
			if err != nil {
				log.Fatalf("error importing csv: %s\n", err)
			}
		}
	}
	w.Header().Add("Content-Type", "text/html")
	w.Write([]byte("<!DOCTYPE html><html><body>Upload successful!</body></html>"))
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
	w.Header().Add("Content-Type", "text/html")
	err = tmpl.Execute(w, State{
		Datasets: datasets,
	})
	if err != nil {
		log.Fatal(err)
	}
}
