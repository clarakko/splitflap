package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"splitflap-api-go/internal/middleware"
	"splitflap-api-go/internal/model"
	"splitflap-api-go/internal/repository"
	"splitflap-api-go/internal/service"

	_ "github.com/mattn/go-sqlite3"
)

const testSchemaSQL = `-- SplitFlap Display Database Schema
-- SQLite version

CREATE TABLE IF NOT EXISTS displays (
    id TEXT PRIMARY KEY,
    content_rows TEXT NOT NULL,  -- JSON-serialized 2D array of strings
    row_count INTEGER NOT NULL CHECK(row_count >= 1 AND row_count <= 20),
    column_count INTEGER NOT NULL CHECK(column_count >= 1 AND column_count <= 10),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index for listing displays (ordered by creation time)
CREATE INDEX IF NOT EXISTS idx_displays_created_at ON displays(created_at);
`

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	if _, err := db.Exec(testSchemaSQL); err != nil {
		db.Close()
		t.Fatalf("failed to execute schema: %v", err)
	}

	return db
}

func newTestHandler(t *testing.T) (http.Handler, repository.DisplayRepository) {
	t.Helper()

	db := setupTestDB(t)
	t.Cleanup(func() { db.Close() })

	repo := repository.NewSqliteDisplayRepository(db)
	svc := service.NewDisplayService(repo)
	return middleware.WithCORS(NewDisplayHandler(svc)), repo
}

func createTestDisplay(id string) *model.Display {
	return &model.Display{
		ID: id,
		Content: model.DisplayContent{
			Rows: [][]string{
				{"H", "E", "L", "L", "O"},
				{"W", "O", "R", "L", "D"},
			},
		},
		Config: model.DisplayConfig{
			RowCount:    2,
			ColumnCount: 5,
		},
	}
}

func TestGetDisplayOK(t *testing.T) {
	handler, repo := newTestHandler(t)

	// Create a test display
	display := createTestDisplay("test-1")
	if err := repo.Create(display); err != nil {
		t.Fatalf("failed to create test display: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/displays/test-1", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	if contentType := rec.Header().Get("Content-Type"); contentType != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", contentType)
	}

	var retrieved model.Display
	if err := json.Unmarshal(rec.Body.Bytes(), &retrieved); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if retrieved.ID != "test-1" {
		t.Fatalf("expected id test-1, got %q", retrieved.ID)
	}

	if retrieved.Config.RowCount != 2 || retrieved.Config.ColumnCount != 5 {
		t.Fatalf("unexpected config: %+v", retrieved.Config)
	}

	if len(retrieved.Content.Rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(retrieved.Content.Rows))
	}

	if retrieved.Content.Rows[0][0] != "H" {
		t.Fatalf("expected first cell to be 'H', got %q", retrieved.Content.Rows[0][0])
	}
}

func TestGetDisplayNotFound(t *testing.T) {
	handler, _ := newTestHandler(t)
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

func TestListDisplays(t *testing.T) {
	handler, repo := newTestHandler(t)

	// Create multiple displays
	display1 := createTestDisplay("test-1")
	display2 := createTestDisplay("test-2")

	if err := repo.Create(display1); err != nil {
		t.Fatalf("failed to create display1: %v", err)
	}
	if err := repo.Create(display2); err != nil {
		t.Fatalf("failed to create display2: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/displays/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var displays []*model.Display
	if err := json.Unmarshal(rec.Body.Bytes(), &displays); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(displays) != 2 {
		t.Fatalf("expected 2 displays, got %d", len(displays))
	}
}

func TestListDisplays_NoTrailingSlash(t *testing.T) {
	handler, repo := newTestHandler(t)

	// Create a display to list
	display := createTestDisplay("test-1")
	if err := repo.Create(display); err != nil {
		t.Fatalf("failed to create display: %v", err)
	}

	// GET without trailing slash (the actual case used by frontend)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/displays", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var displays []*model.Display
	if err := json.Unmarshal(rec.Body.Bytes(), &displays); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(displays) != 1 {
		t.Fatalf("expected 1 display, got %d", len(displays))
	}
}

func TestCreateDisplay(t *testing.T) {
	handler, _ := newTestHandler(t)

	display := createTestDisplay("new-display")
	body, _ := json.Marshal(display)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/displays/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: %s", rec.Code, rec.Body.String())
	}

	var created model.Display
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if created.ID != "new-display" {
		t.Fatalf("expected id new-display, got %q", created.ID)
	}
}

func TestCreateDisplay_NoTrailingSlash(t *testing.T) {
	handler, repo := newTestHandler(t)

	display := createTestDisplay("new-display-no-slash")
	body, _ := json.Marshal(display)

	// POST without trailing slash - this was returning 301 before the fix
	req := httptest.NewRequest(http.MethodPost, "/api/v1/displays", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: %s", rec.Code, rec.Body.String())
	}

	var created model.Display
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if created.ID != "new-display-no-slash" {
		t.Fatalf("expected id new-display-no-slash, got %q", created.ID)
	}

	// Verify it was actually persisted to the database
	retrieved, err := repo.GetByID("new-display-no-slash")
	if err != nil {
		t.Fatalf("failed to retrieve created display from database: %v", err)
	}
	if retrieved == nil {
		t.Fatalf("expected display to be persisted, but got nil")
	}
}

func TestCreateDisplay_InvalidRequest(t *testing.T) {
	handler, _ := newTestHandler(t)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/displays/", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestCreateDisplay_ValidationError(t *testing.T) {
	handler, _ := newTestHandler(t)

	// Display with invalid rowCount (0)
	display := createTestDisplay("invalid")
	display.Config.RowCount = 0

	body, _ := json.Marshal(display)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/displays/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestCreateDisplay_Duplicate(t *testing.T) {
	handler, repo := newTestHandler(t)

	display := createTestDisplay("duplicate")
	if err := repo.Create(display); err != nil {
		t.Fatalf("failed to create first display: %v", err)
	}

	body, _ := json.Marshal(display)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/displays/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("expected status 409, got %d", rec.Code)
	}
}

func TestUpdateDisplay(t *testing.T) {
	handler, repo := newTestHandler(t)

	display := createTestDisplay("test-1")
	if err := repo.Create(display); err != nil {
		t.Fatalf("failed to create display: %v", err)
	}

	// Update the display
	display.Content.Rows[0][0] = "X"
	body, _ := json.Marshal(display)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/displays/test-1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var updated model.Display
	if err := json.Unmarshal(rec.Body.Bytes(), &updated); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if updated.Content.Rows[0][0] != "X" {
		t.Fatalf("expected first cell to be 'X', got %q", updated.Content.Rows[0][0])
	}
}

func TestUpdateDisplay_NotFound(t *testing.T) {
	handler, _ := newTestHandler(t)

	display := createTestDisplay("nonexistent")
	body, _ := json.Marshal(display)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/displays/nonexistent", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}
}

func TestUpdateDisplay_IDMismatch(t *testing.T) {
	handler, _ := newTestHandler(t)

	display := createTestDisplay("test-1")
	body, _ := json.Marshal(display)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/displays/different-id", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestUpdateDisplay_NoTrailingSlash(t *testing.T) {
	handler, repo := newTestHandler(t)

	display := createTestDisplay("test-update-no-slash")
	if err := repo.Create(display); err != nil {
		t.Fatalf("failed to create display: %v", err)
	}

	// Update the display
	display.Content.Rows[0][0] = "X"
	body, _ := json.Marshal(display)

	// PUT without trailing slash - this was failing with 404 before the CORS fix
	req := httptest.NewRequest(http.MethodPut, "/api/v1/displays/test-update-no-slash", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var updated model.Display
	if err := json.Unmarshal(rec.Body.Bytes(), &updated); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if updated.Content.Rows[0][0] != "X" {
		t.Fatalf("expected first cell to be 'X', got %q", updated.Content.Rows[0][0])
	}

	// Verify it was actually persisted to the database
	retrieved, err := repo.GetByID("test-update-no-slash")
	if err != nil {
		t.Fatalf("failed to retrieve updated display from database: %v", err)
	}
	if retrieved == nil {
		t.Fatalf("expected display to be persisted, but got nil")
	}
	if retrieved.Content.Rows[0][0] != "X" {
		t.Fatalf("expected persisted display's first cell to be 'X', got %q", retrieved.Content.Rows[0][0])
	}
}

func TestDeleteDisplay(t *testing.T) {
	handler, repo := newTestHandler(t)

	display := createTestDisplay("test-1")
	if err := repo.Create(display); err != nil {
		t.Fatalf("failed to create display: %v", err)
	}

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/displays/test-1", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", rec.Code)
	}

	// Verify it was deleted
	_, err := repo.GetByID("test-1")
	if err != repository.ErrNotFound {
		t.Fatalf("expected display to be deleted, but got: %v", err)
	}
}

func TestDeleteDisplay_NotFound(t *testing.T) {
	handler, _ := newTestHandler(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/displays/nonexistent", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}
}

func TestDeleteDisplay_NoTrailingSlash(t *testing.T) {
	handler, repo := newTestHandler(t)

	display := createTestDisplay("test-delete-no-slash")
	if err := repo.Create(display); err != nil {
		t.Fatalf("failed to create display: %v", err)
	}

	// DELETE without trailing slash - should work correctly
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/displays/test-delete-no-slash", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", rec.Code)
	}

	// Verify it was deleted
	_, err := repo.GetByID("test-delete-no-slash")
	if err != repository.ErrNotFound {
		t.Fatalf("expected display to be deleted, but got: %v", err)
	}
}

func TestCORSAllowedOrigin(t *testing.T) {
	handler, _ := newTestHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/displays/", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if origin := rec.Header().Get("Access-Control-Allow-Origin"); origin != "http://localhost:5173" {
		t.Fatalf("expected Access-Control-Allow-Origin to echo origin, got %q", origin)
	}
}

func TestCORSAllowedMethods(t *testing.T) {
	handler, _ := newTestHandler(t)
	req := httptest.NewRequest(http.MethodOptions, "/api/v1/displays", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", rec.Code)
	}

	methods := rec.Header().Get("Access-Control-Allow-Methods")
	if methods == "" {
		t.Fatalf("expected Access-Control-Allow-Methods header, got empty")
	}

	// Verify all required methods are allowed
	requiredMethods := []string{"GET", "POST", "PUT", "DELETE"}
	for _, method := range requiredMethods {
		if !contains(methods, method) {
			t.Fatalf("expected %s in Access-Control-Allow-Methods, got %q", method, methods)
		}
	}
}

func contains(s, substr string) bool {
	// Simple substring check for methods in comma-separated list
	start := 0
	for {
		idx := strings.Index(s[start:], substr)
		if idx == -1 {
			return false
		}
		idx += start
		// Check it's a complete word (surrounded by non-alphanumeric or boundaries)
		if (idx == 0 || !isAlpha(rune(s[idx-1]))) && (idx+len(substr) >= len(s) || !isAlpha(rune(s[idx+len(substr)]))) {
			return true
		}
		start = idx + 1
	}
}

func isAlpha(r rune) bool {
	return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z')
}
