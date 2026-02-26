package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"splitflap-api-go/internal/model"
	"splitflap-api-go/internal/repository"
	"splitflap-api-go/internal/service"
)

const displayBasePath = "/api/v1/displays/"

type DisplayHandler struct {
	service *service.DisplayService
}

func NewDisplayHandler(service *service.DisplayService) http.Handler {
	return &DisplayHandler{service: service}
}

func (h *DisplayHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, displayBasePath) {
		http.NotFound(w, r)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, displayBasePath)

	// Route based on method and whether ID is present
	switch r.Method {
	case http.MethodGet:
		if id == "" {
			h.handleList(w, r)
		} else if strings.Contains(id, "/") {
			http.NotFound(w, r)
		} else {
			h.handleGet(w, r, id)
		}
	case http.MethodPost:
		if id != "" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		h.handleCreate(w, r)
	case http.MethodPut:
		if id == "" || strings.Contains(id, "/") {
			http.NotFound(w, r)
			return
		}
		h.handleUpdate(w, r, id)
	case http.MethodDelete:
		if id == "" || strings.Contains(id, "/") {
			http.NotFound(w, r)
			return
		}
		h.handleDelete(w, r, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *DisplayHandler) handleGet(w http.ResponseWriter, r *http.Request, id string) {
	display := h.service.GetDisplay(id)
	if display == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	h.respondJSON(w, http.StatusOK, display)
}

func (h *DisplayHandler) handleList(w http.ResponseWriter, r *http.Request) {
	displays, err := h.service.ListDisplays()
	if err != nil {
		log.Printf("failed to list displays: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, http.StatusOK, displays)
}

func (h *DisplayHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var display model.Display
	if err := json.NewDecoder(r.Body).Decode(&display); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	// Validate display
	if err := validateDisplay(&display); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if err := h.service.CreateDisplay(&display); err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{"error": "display already exists"})
			return
		}
		log.Printf("failed to create display: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, http.StatusCreated, display)
}

func (h *DisplayHandler) handleUpdate(w http.ResponseWriter, r *http.Request, id string) {
	var display model.Display
	if err := json.NewDecoder(r.Body).Decode(&display); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	// Ensure ID in URL matches ID in body
	if display.ID != id {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "ID mismatch"})
		return
	}

	// Validate display
	if err := validateDisplay(&display); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if err := h.service.UpdateDisplay(&display); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Printf("failed to update display: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, http.StatusOK, display)
}

func (h *DisplayHandler) handleDelete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.service.DeleteDisplay(id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Printf("failed to delete display: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *DisplayHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("failed to encode JSON response: %v", err)
	}
}

func validateDisplay(display *model.Display) error {
	if display.ID == "" {
		return errors.New("display ID is required")
	}

	// Validate row and column counts per DATA_MODEL.md constraints
	if display.Config.RowCount < 1 || display.Config.RowCount > 20 {
		return errors.New("rowCount must be between 1 and 20")
	}
	if display.Config.ColumnCount < 1 || display.Config.ColumnCount > 10 {
		return errors.New("columnCount must be between 1 and 10")
	}

	// Validate content matches config
	if len(display.Content.Rows) != display.Config.RowCount {
		return errors.New("content rows count must match rowCount")
	}
	for i, row := range display.Content.Rows {
		if len(row) != display.Config.ColumnCount {
			return errors.New("row " + string(rune(i)) + " column count must match columnCount")
		}
		// Validate each cell is 0 or 1 characters (empty or single character)
		for j, cell := range row {
			if len(cell) > 1 {
				return errors.New("cell at row " + string(rune(i)) + " column " + string(rune(j)) + " must be 0 or 1 characters")
			}
		}
	}

	return nil
}
