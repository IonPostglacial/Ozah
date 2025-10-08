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
	"time"

	"embed"

	"nicolas.galipot.net/hazo/server/appdb"
	"nicolas.galipot.net/hazo/server/authentication"
	"nicolas.galipot.net/hazo/server/common"
	"nicolas.galipot.net/hazo/server/components"
	"nicolas.galipot.net/hazo/server/documents"
	"nicolas.galipot.net/hazo/server/export"
	"nicolas.galipot.net/hazo/server/views/characters"
	"nicolas.galipot.net/hazo/server/views/identification"
	"nicolas.galipot.net/hazo/server/views/index"
	"nicolas.galipot.net/hazo/server/views/taxons"
	"nicolas.galipot.net/hazo/storage/dataset"
)

//go:embed assets
var assets embed.FS

//go:embed debug.js
var debugJS string

type Server struct {
	*http.ServeMux
}

func New(config *common.ServerConfig) Server {
	s := http.NewServeMux()
	s.HandleFunc("/ds/{dsName}/json", common.Handler(export.JsonHandler).
		Wrap(authentication.HandlerWrapper).
		Wrap(appdb.Handler).
		Unwrap(config))
	s.HandleFunc("/ds/{dsName}/taxons", common.Handler(taxons.Handler).
		Wrap(authentication.HandlerWrapper).
		Wrap(documents.HandlerWrapper("taxons")).
		Wrap(appdb.Handler).
		Unwrap(config))
	s.HandleFunc("/ds/{dsName}/taxons/{id}", common.Handler(taxons.Handler).
		Wrap(authentication.HandlerWrapper).
		Wrap(documents.HandlerWrapper("taxons")).
		Wrap(appdb.Handler).
		Unwrap(config))
	s.HandleFunc("/ds/{dsName}/characters", common.Handler(characters.Handler).
		Wrap(authentication.HandlerWrapper).
		Wrap(documents.HandlerWrapper("characters")).
		Wrap(appdb.Handler).
		Unwrap(config))
	s.HandleFunc("/ds/{dsName}/characters/{id}", common.Handler(characters.Handler).
		Wrap(authentication.HandlerWrapper).
		Wrap(documents.HandlerWrapper("characters")).
		Wrap(appdb.Handler).
		Unwrap(config))
	s.HandleFunc("/ds/{dsName}/identify", common.Handler(identification.Handler).
		Wrap(authentication.HandlerWrapper).
		Wrap(appdb.Handler).
		Unwrap(config))
	s.HandleFunc("/upload", common.Handler(uploadHandler).
		Wrap(authentication.HandlerWrapper).
		Wrap(appdb.Handler).
		Unwrap(config))
	s.HandleFunc("/", common.Handler(index.Handler).
		Wrap(authentication.HandlerWrapper).
		Wrap(appdb.Handler).
		Unwrap(config))
	s.HandleFunc("/components.js", common.Handler(components.JavascriptHandler).Unwrap(config))
	s.HandleFunc("/components.css", common.Handler(components.CssHandler).Unwrap(config))
	s.Handle("/assets/", http.FileServer(http.FS(assets)))
	if config.Debug {
		startedOn := time.Now().Format(time.RFC3339)

		s.HandleFunc("/started-on", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, startedOn)
		})
		s.HandleFunc("/debug.js", func(w http.ResponseWriter, r *http.Request) {
			io.Copy(w, strings.NewReader(debugJS))
		})
	}
	return Server{s}
}

func createFile(p string) (*os.File, error) {
	if err := os.MkdirAll(path.Dir(p), 0770); err != nil {
		return nil, err
	}
	return os.Create(p)
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
			return fmt.Errorf("failed to create temporary directory to upload file '%s': %w", fileName, err)
		}
		defer os.RemoveAll(dir)
		file, _ := io.ReadAll(auxiliar)
		isZip := strings.HasSuffix(fileName, ".zip")
		dbName := ""
		if isZip {
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
				file, err := createFile(filePath)
				if err != nil {
					return fmt.Errorf("creating file '%s' from '%s', '%s' failed: %w", filePath, dir, f.Name, err)
				}
				defer file.Close()
				_, err = io.Copy(file, content)
				if err != nil {
					return fmt.Errorf("copying file '%s' failed: %w", filePath, err)
				}
			}
			dbName = strings.TrimSuffix(fileName, ".zip")
		} else {
			dbName = strings.TrimSuffix(fileName, ".hazo.json")
		}
		dbPath, err := cc.User.GetDataset(dbName)
		if err != nil {
			return fmt.Errorf("uploading file '%s' failed while retrieving user dataset: %w", fileName, err)
		}
		err = dataset.Create(dbPath)
		if err != nil {
			return fmt.Errorf("creating database '%s' failed: %w", dbPath, err)
		}
		if isZip {
			err = dataset.ImportCsv(dir, dbPath)
		} else {
			err = dataset.ImportJson(file, dbPath)
		}
		if err != nil {
			return fmt.Errorf("error importing dataset '%s' to '%s' failed: %w", dir, dbPath, err)
		}
	}
	w.Header().Add("Content-Type", "text/html")
	w.Write([]byte("<!DOCTYPE html><html><body><div class='upload-msg'>Upload successful!</div></body></html>"))
	return nil
}
