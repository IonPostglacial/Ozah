package server

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"embed"

	"nicolas.galipot.net/hazo/db"
	"nicolas.galipot.net/hazo/server/authentication"
	"nicolas.galipot.net/hazo/server/common"
	"nicolas.galipot.net/hazo/server/components"
	"nicolas.galipot.net/hazo/server/documents"
	"nicolas.galipot.net/hazo/server/views/characters"
	"nicolas.galipot.net/hazo/server/views/identification"
	"nicolas.galipot.net/hazo/server/views/taxons"
)

//go:embed assets
var assets embed.FS

//go:embed index.html
var indexPage string

func New() *http.ServeMux {
	s := http.NewServeMux()
	s.HandleFunc("/ds/{dsName}/taxons", common.Handler(taxons.Handler).
		Wrap(authentication.HandlerWrapper).
		Wrap(documents.HandlerWrapper("taxons")).
		Unwrap())
	s.HandleFunc("/ds/{dsName}/taxons/{id}", common.Handler(taxons.Handler).
		Wrap(authentication.HandlerWrapper).
		Wrap(documents.HandlerWrapper("taxons")).
		Unwrap())
	s.HandleFunc("/ds/{dsName}/characters", common.Handler(characters.Handler).
		Wrap(authentication.HandlerWrapper).
		Wrap(documents.HandlerWrapper("characters")).
		Unwrap())
	s.HandleFunc("/ds/{dsName}/characters/{id}", common.Handler(characters.Handler).
		Wrap(authentication.HandlerWrapper).
		Wrap(documents.HandlerWrapper("characters")).
		Unwrap())
	s.HandleFunc("/ds/{dsName}/identify", common.Handler(identification.Handler).
		Wrap(authentication.HandlerWrapper).
		Unwrap())
	s.HandleFunc("/upload", common.Handler(uploadHandler).
		Wrap(authentication.HandlerWrapper).
		Unwrap())
	s.HandleFunc("/", common.Handler(indexHandler).
		Wrap(authentication.HandlerWrapper).
		Unwrap())
	s.HandleFunc("/components.js", common.Handler(components.JavascriptHandler).Unwrap())
	s.HandleFunc("/components.css", common.Handler(components.CssHandler).Unwrap())
	s.Handle("/assets/", http.FileServer(http.FS(assets)))
	return s
}

type State struct {
	PageTitle string
	Datasets  []db.Dataset
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
			return fmt.Errorf("reading zip failed: %w", err)
		}
		for _, f := range r.File {
			content, err := f.Open()
			if err != nil {
				return fmt.Errorf("reading file '%s' failed: %w", f.Name, err)
			}
			filePath := path.Join(dir, f.Name)
			file, err := os.Create(filePath)
			if err != nil {
				return fmt.Errorf("creating file '%s' failed: %w", filePath, err)
			}
			defer file.Close()
			_, err = io.Copy(file, content)
			if err != nil {
				return fmt.Errorf("copying file '%s' failed: %w", filePath, err)
			}
		}
		dbName := strings.TrimSuffix(fileName, ".zip")
		dbPath, err := cc.User.GetDataset(dbName)
		if err != nil {
			return err
		}
		err = db.Create(dbPath)
		if err != nil {
			return fmt.Errorf("creating database '%s' failed: %w", dbPath, err)
		}
		err = db.ImportCsv(dir, dbPath)
		if err != nil {
			return fmt.Errorf("error importing zip '%s' to '%s' failed: %w", dir, dbPath, err)
		}
	}
	w.Header().Add("Content-Type", "text/html")
	w.Write([]byte("<!DOCTYPE html><html><body><div class='upload-msg'>Upload successful!</div></body></html>"))
	return nil
}

func indexHandler(w http.ResponseWriter, r *http.Request, cc *common.Context) error {
	tmpl := components.NewTemplate()
	tmpl, err := tmpl.Parse(indexPage)
	if err != nil {
		return err
	}
	datasets, err := cc.User.ListDatasets()
	if err != nil {
		return err
	}
	w.Header().Add("Content-Type", "text/html")
	err = tmpl.Execute(w, State{
		PageTitle: "Hazo Home",
		Datasets:  datasets,
	})
	if err != nil {
		return err
	}
	return nil
}
