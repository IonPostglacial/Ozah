package server

import (
	"archive/zip"
	"bytes"
	"fmt"
	"html/template"
	"io"
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
	s.HandleFunc("/ds/{dsName}/taxons", common.UnwrapHandler(authentication.HandlerWrapper(taxons.Handler)))
	s.HandleFunc("/ds/{dsName}/taxons/{id}", common.UnwrapHandler(authentication.HandlerWrapper(taxons.Handler)))
	s.HandleFunc("/upload", common.UnwrapHandler(authentication.HandlerWrapper(uploadHandler)))
	s.HandleFunc("/", common.UnwrapHandler(authentication.HandlerWrapper(indexHandler)))
	s.Handle("/assets/", http.FileServer(http.FS(assets)))
	return s
}

type State struct {
	Datasets []db.Dataset
}

func uploadHandler(w http.ResponseWriter, r *http.Request, cc *common.Context) error {
	r.ParseMultipartForm(5000000) //5 MB in memory, the rest in disk
	datas := r.MultipartForm
	for _, headers := range datas.File {
		if len(headers) != 1 {
			return fmt.Errorf("wrong header length, expected 1 got %d", len(headers))
		}
		auxiliar, _ := headers[0].Open()
		fileName := headers[0].Filename
		dir, err := os.MkdirTemp("tmp", fileName)
		if err != nil {
			return err
		}
		defer os.RemoveAll(dir)
		file, _ := io.ReadAll(auxiliar)
		r, err := zip.NewReader(bytes.NewReader(file), int64(len(file)))
		if err != nil {
			return fmt.Errorf("error reading zip: %w", err)
		}
		for _, f := range r.File {
			content, err := f.Open()
			if err != nil {
				return fmt.Errorf("error reading file '%s': %w", f.Name, err)
			}
			filePath := path.Join(dir, f.Name)
			file, err := os.Create(filePath)
			if err != nil {
				return fmt.Errorf("error creating file '%s': %w", filePath, err)
			}
			defer file.Close()
			_, err = io.Copy(file, content)
			if err != nil {
				return fmt.Errorf("error copying file '%s': %w", filePath, err)
			}
		}
		dbPath := fmt.Sprintf("%s.sq3", strings.TrimSuffix(fileName, ".zip"))
		err = db.Init(dbPath)
		if err != nil {
			return fmt.Errorf("error creating db '%s': %w", dbPath, err)
		}
		err = db.ImportCsv(dir, dbPath)
		if err != nil {
			return fmt.Errorf("error importing zip: %w", err)
		}
	}
	w.Header().Add("Content-Type", "text/html")
	w.Write([]byte("<!DOCTYPE html><html><body>Upload successful!</body></html>"))
	return nil
}

func indexHandler(w http.ResponseWriter, r *http.Request, cc *common.Context) error {
	tmpl := template.New("index")
	tmpl, err := tmpl.Parse(indexPage)
	if err != nil {
		return err
	}
	datasets, err := db.ListDatasets()
	if err != nil {
		return err
	}
	w.Header().Add("Content-Type", "text/html")
	err = tmpl.Execute(w, State{
		Datasets: datasets,
	})
	if err != nil {
		return err
	}
	return nil
}
