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

type DatasetExportResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

func DatasetExportJsonHandler(w http.ResponseWriter, r *http.Request, cc *common.Context) error {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(DatasetExportResponse{
			Success: false,
			Error:   "Method not allowed. Use GET.",
		})
		return nil
	}

	dsName := r.PathValue("name")
	if dsName == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(DatasetExportResponse{
			Success: false,
			Error:   "Missing dataset name in URL path",
		})
		return nil
	}

	pds, err := cc.User.CanAccessDataset(dsName)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(DatasetExportResponse{
			Success: false,
			Error:   fmt.Sprintf("Access denied: %v", err),
		})
		return nil
	}

	queries, err := dataset.OpenDb(pds)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(DatasetExportResponse{
			Success: false,
			Error:   fmt.Sprintf("Could not open dataset database: %v", err),
		})
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.hazo.json", dsName))
	w.WriteHeader(http.StatusOK)

	if err := dataset.ExportJson(dsName, queries, w); err != nil {
		return fmt.Errorf("failed to export dataset '%s' as JSON: %w", dsName, err)
	}

	return nil
}

func DatasetExportCsvHandler(w http.ResponseWriter, r *http.Request, cc *common.Context) error {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(DatasetExportResponse{
			Success: false,
			Error:   "Method not allowed. Use GET.",
		})
		return nil
	}

	dsName := r.PathValue("name")
	if dsName == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(DatasetExportResponse{
			Success: false,
			Error:   "Missing dataset name in URL path",
		})
		return nil
	}

	pds, err := cc.User.CanAccessDataset(dsName)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(DatasetExportResponse{
			Success: false,
			Error:   fmt.Sprintf("Access denied: %v", err),
		})
		return nil
	}

	queries, err := dataset.OpenDb(pds)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(DatasetExportResponse{
			Success: false,
			Error:   fmt.Sprintf("Could not open dataset database: %v", err),
		})
		return nil
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.zip", dsName))
	w.WriteHeader(http.StatusOK)

	if err := dataset.ExportCsv(dsName, queries, w); err != nil {
		return fmt.Errorf("failed to export dataset '%s' as CSV: %w", dsName, err)
	}

	return nil
}
