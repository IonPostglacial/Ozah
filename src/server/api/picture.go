package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"nicolas.galipot.net/hazo/server/common"
	"nicolas.galipot.net/hazo/server/picture"
)

type PictureUploadResponse struct {
	Success         bool   `json:"success"`
	Message         string `json:"message,omitempty"`
	AttachmentIndex int    `json:"attachmentIndex,omitempty"`
	Error           string `json:"error,omitempty"`
}

func PictureUploadHandler(w http.ResponseWriter, r *http.Request, cc *common.Context) error {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(PictureUploadResponse{
			Success: false,
			Error:   "Method not allowed. Use POST.",
		})
		return nil
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil { // 32 MB max
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(PictureUploadResponse{
			Success: false,
			Error:   fmt.Sprintf("Could not parse multipart form: %v", err),
		})
		return nil
	}

	dsName := r.FormValue("dataset")
	docRef := r.FormValue("ref")
	indexStr := r.FormValue("index")

	if dsName == "" || docRef == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(PictureUploadResponse{
			Success: false,
			Error:   "Missing required fields: dataset, ref",
		})
		return nil
	}

	attachmentIndex := -1
	if indexStr != "" {
		var err error
		attachmentIndex, err = strconv.Atoi(indexStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(PictureUploadResponse{
				Success: false,
				Error:   "Invalid attachment index",
			})
			return nil
		}
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(PictureUploadResponse{
			Success: false,
			Error:   fmt.Sprintf("Could not get file: %v", err),
		})
		return nil
	}
	defer file.Close()

	ctx := context.Background()
	result, err := picture.UploadPictureFromMultipart(ctx, cc.User, dsName, docRef, file, header, attachmentIndex)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(PictureUploadResponse{
			Success: false,
			Error:   err.Error(),
		})
		return nil
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(PictureUploadResponse{
		Success:         true,
		Message:         "Picture uploaded successfully",
		AttachmentIndex: result.AttachmentIndex,
	})
	return nil
}
