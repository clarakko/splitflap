package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

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
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if !strings.HasPrefix(r.URL.Path, displayBasePath) {
		http.NotFound(w, r)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, displayBasePath)
	if id == "" || strings.Contains(id, "/") {
		http.NotFound(w, r)
		return
	}

	display := h.service.GetDisplay(id)
	if display == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(display); err != nil {
		log.Printf("failed to encode display response: %v", err)
	}
}
