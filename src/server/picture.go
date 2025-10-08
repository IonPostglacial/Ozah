package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"nicolas.galipot.net/hazo/server/common"
	"nicolas.galipot.net/hazo/storage/dataset"
	"nicolas.galipot.net/hazo/storage/dsdb"
)

func PictureHandler(w http.ResponseWriter, r *http.Request, cc *common.Context) error {
	dsName := r.PathValue("dsName")
	docRef := r.PathValue("docRef")
	indexStr := r.PathValue("index")

	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return fmt.Errorf("invalid attachment index: %w", err)
	}

	size := dataset.ParseThumbnailSize(r.URL.Query().Get("size"))

	ds, err := cc.User.GetDataset(dsName)
	if err != nil {
		sharedDs, err := cc.User.GetReadableSharedDatasets()
		if err != nil {
			return fmt.Errorf("could not access dataset: %w", err)
		}

		found := false
		for _, shared := range sharedDs {
			if shared.Name == dsName {
				ds = dataset.Private(shared.Path)
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("access denied to dataset '%s'", dsName)
		}
	}

	queries, err := dataset.OpenDb(ds)
	if err != nil {
		return fmt.Errorf("could not open dataset: %w", err)
	}

	ctx := context.Background()

	attachment, err := queries.GetDocumentAttachmentByIndex(ctx, dsdb.GetDocumentAttachmentByIndexParams{
		DocumentRef:     docRef,
		AttachmentIndex: int64(index),
	})
	if err != nil {
		return fmt.Errorf("attachment not found: %w", err)
	}

	var filePath string
	switch size {
	case dataset.SizeSmall:
		filePath = attachment.PathSmall
	case dataset.SizeMedium:
		filePath = attachment.PathMedium
	case dataset.SizeBig:
		filePath = attachment.PathBig
	default:
		filePath = attachment.PathMedium
	}

	if filePath == "" || !fileExists(filePath) {
		filePath = attachment.Path
	}

	if !fileExists(filePath) {
		return fmt.Errorf("picture file not found: %s", filePath)
	}

	contentType := "image/jpeg"
	ext := filepath.Ext(filePath)
	switch ext {
	case ".png":
		contentType = "image/png"
	case ".gif":
		contentType = "image/gif"
	case ".webp":
		contentType = "image/webp"
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Cache-Control", "public, max-age=86400") // Cache for 24 hours

	http.ServeFile(w, r, filePath)
	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
