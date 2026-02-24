package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"splitflap-api-go/internal/middleware"
	"splitflap-api-go/internal/model"
	"splitflap-api-go/internal/service"
)

func newTestHandler() http.Handler {
	svc := service.NewDisplayService()
	return middleware.WithCORS(NewDisplayHandler(svc))
}

func TestGetDisplayOK(t *testing.T) {
	handler := newTestHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/displays/demo", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	if contentType := rec.Header().Get("Content-Type"); contentType != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", contentType)
	}

	var display model.Display
	if err := json.Unmarshal(rec.Body.Bytes(), &display); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if display.ID != "demo" {
		t.Fatalf("expected id demo, got %q", display.ID)
	}

	if display.Config.RowCount != 5 || display.Config.ColumnCount != 4 {
		t.Fatalf("unexpected config: %+v", display.Config)
	}

	if len(display.Content.Rows) != 5 {
		t.Fatalf("expected 5 rows, got %d", len(display.Content.Rows))
	}

	headerRow := display.Content.Rows[0]
	expectedHeader := []string{"TIME", "DESTINATION", "PLATFORM", "STATUS"}
	for i, value := range expectedHeader {
		if headerRow[i] != value {
			t.Fatalf("expected header %q at index %d, got %q", value, i, headerRow[i])
		}
	}
}

func TestGetDisplayNotFound(t *testing.T) {
	handler := newTestHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/displays/nonexistent", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}

	if body := rec.Body.String(); body != "" {
		t.Fatalf("expected empty body, got %q", body)
	}
}

func TestCORSAllowedOrigin(t *testing.T) {
	handler := newTestHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/displays/demo", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if origin := rec.Header().Get("Access-Control-Allow-Origin"); origin != "http://localhost:5173" {
		t.Fatalf("expected Access-Control-Allow-Origin to echo origin, got %q", origin)
	}
}
