package server

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"embed"

	"nicolas.galipot.net/hazo/server/api"
	"nicolas.galipot.net/hazo/server/appdb"
	"nicolas.galipot.net/hazo/server/authentication"
	"nicolas.galipot.net/hazo/server/common"
	"nicolas.galipot.net/hazo/server/components"
	"nicolas.galipot.net/hazo/server/documents"
	"nicolas.galipot.net/hazo/server/export"
	"nicolas.galipot.net/hazo/server/views/admin"
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
	s.HandleFunc("/auth/microsoft/login", common.Handler(authentication.MSLoginHandler).
		Wrap(appdb.Handler).
		Unwrap(config))
	s.HandleFunc("/auth/microsoft/callback", common.Handler(authentication.MSCallbackHandler).
		Wrap(appdb.Handler).
		Unwrap(config))
	s.HandleFunc("/logout", common.Handler(authentication.LogoutHandler).
		Wrap(appdb.Handler).
		Unwrap(config))
	s.HandleFunc("/admin", common.Handler(admin.Handler).
		Wrap(authentication.RequireCapability("user.manage")).
		Wrap(authentication.HandlerWrapper).
		Wrap(appdb.Handler).
		Unwrap(config))
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
	s.HandleFunc("/ds/{dsName}/pictures/{docRef}/{index}", common.Handler(PictureHandler).
		Wrap(authentication.HandlerWrapper).
		Wrap(appdb.Handler).
		Unwrap(config))
	s.HandleFunc("/api/pictures", common.Handler(api.PictureUploadHandler).
		Wrap(authentication.HandlerWrapper).
		Wrap(appdb.Handler).
		Unwrap(config))
	s.HandleFunc("/api/datasets/json", common.Handler(api.DatasetImportJsonHandler).
		Wrap(authentication.HandlerWrapper).
		Wrap(appdb.Handler).
		Unwrap(config))
	s.HandleFunc("/api/datasets/csv", common.Handler(api.DatasetImportCsvHandler).
		Wrap(authentication.HandlerWrapper).
		Wrap(appdb.Handler).
		Unwrap(config))
	s.HandleFunc("/api/datasets/{name}/json", common.Handler(api.DatasetExportJsonHandler).
		Wrap(authentication.HandlerWrapper).
		Wrap(appdb.Handler).
		Unwrap(config))
	s.HandleFunc("/api/datasets/{name}/csv", common.Handler(api.DatasetExportCsvHandler).
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

func uploadHandler(w http.ResponseWriter, r *http.Request, cc *common.Context) error {
	r.ParseMultipartForm(5000000) //5 MB in memory, the rest in disk
	datas := r.MultipartForm
	for _, headers := range datas.File {
		if len(headers) != 1 {
			return fmt.Errorf("wrong header length, expected 1 got %d", len(headers))
		}
		fileReader, _ := headers[0].Open()
		fileName := headers[0].Filename
		fileData, _ := io.ReadAll(fileReader)

		isZip := strings.HasSuffix(fileName, ".zip")
		var dbName string
		if isZip {
			dbName = strings.TrimSuffix(fileName, ".zip")
		} else {
			dbName = strings.TrimSuffix(fileName, ".hazo.json")
		}

		dbPath, err := cc.User.GetDataset(dbName)
		if err != nil {
			return fmt.Errorf("uploading file '%s' failed while retrieving user dataset: %w", fileName, err)
		}

		if isZip {
			err = dataset.ImportCsvDataset(dataset.Private(dbPath), fileData)
		} else {
			err = dataset.ImportJsonDataset(dataset.Private(dbPath), fileData)
		}
		if err != nil {
			return fmt.Errorf("error importing dataset from '%s': %w", fileName, err)
		}
	}
	w.Header().Add("Content-Type", "text/html")
	w.Write([]byte("<!DOCTYPE html><html><body><div class='upload-msg'>Upload successful!</div></body></html>"))
	return nil
}
