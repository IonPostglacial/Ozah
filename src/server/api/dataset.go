package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"nicolas.galipot.net/hazo/server/common"
	"nicolas.galipot.net/hazo/storage/dataset"
)

type DatasetImportResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Dataset string `json:"dataset,omitempty"`
	Error   string `json:"error,omitempty"`
}

func validateImportRequest(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(DatasetImportResponse{
			Success: false,
			Error:   "Method not allowed. Use POST.",
		})
		return false
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil { // 32 MB max
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(DatasetImportResponse{
			Success: false,
			Error:   fmt.Sprintf("Could not parse multipart form: %v", err),
		})
		return false
	}

	return true
}

func writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(DatasetImportResponse{
		Success: false,
		Error:   message,
	})
}

func writeSuccessResponse(w http.ResponseWriter, dsName string) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(DatasetImportResponse{
		Success: true,
		Message: "Dataset imported successfully",
		Dataset: dsName,
	})
}

func DatasetImportJsonHandler(w http.ResponseWriter, r *http.Request, cc *common.Context) error {
	w.Header().Set("Content-Type", "application/json")

	if !validateImportRequest(w, r) {
		return nil
	}

	dsName := r.FormValue("dataset")
	if dsName == "" {
		writeErrorResponse(w, http.StatusBadRequest, "Missing required field: dataset")
		return nil
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Could not get file: %v", err))
		return nil
	}
	defer file.Close()

	if !strings.HasSuffix(header.Filename, ".json") && !strings.HasSuffix(header.Filename, ".hazo.json") {
		writeErrorResponse(w, http.StatusBadRequest, "File must be a JSON file (.json or .hazo.json)")
		return nil
	}

	fileData, err := io.ReadAll(file)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Could not read file: %v", err))
		return nil
	}

	dbPath, err := cc.User.GetDataset(dsName)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Could not get dataset path: %v", err))
		return nil
	}

	if err := dataset.ImportJsonDataset(dataset.Private(dbPath), fileData); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Could not import dataset: %v", err))
		return nil
	}

	writeSuccessResponse(w, dsName)
	return nil
}

func DatasetImportCsvHandler(w http.ResponseWriter, r *http.Request, cc *common.Context) error {
	w.Header().Set("Content-Type", "application/json")

	if !validateImportRequest(w, r) {
		return nil
	}

	dsName := r.FormValue("dataset")
	if dsName == "" {
		writeErrorResponse(w, http.StatusBadRequest, "Missing required field: dataset")
		return nil
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Could not get file: %v", err))
		return nil
	}
	defer file.Close()

	if !strings.HasSuffix(header.Filename, ".zip") {
		writeErrorResponse(w, http.StatusBadRequest, "File must be a ZIP archive containing CSV files")
		return nil
	}

	fileData, err := io.ReadAll(file)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Could not read file: %v", err))
		return nil
	}

	dbPath, err := cc.User.GetDataset(dsName)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Could not get dataset path: %v", err))
		return nil
	}

	if err := dataset.ImportCsvDataset(dataset.Private(dbPath), fileData); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Could not import dataset: %v", err))
		return nil
	}

	writeSuccessResponse(w, dsName)
	return nil
}
