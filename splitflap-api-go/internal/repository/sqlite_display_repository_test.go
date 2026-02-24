package repository

import (
	"database/sql"
	_ "embed"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"splitflap-api-go/internal/model"
)

//go:embed schema.sql
var testSchemaSQL string

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

func TestCreate(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSqliteDisplayRepository(db)
	display := createTestDisplay("test-1")

	err := repo.Create(display)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestGetByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSqliteDisplayRepository(db)
	display := createTestDisplay("test-1")

	if err := repo.Create(display); err != nil {
		t.Fatalf("failed to create display: %v", err)
	}

	retrieved, err := repo.GetByID("test-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if retrieved.ID != "test-1" {
		t.Errorf("expected ID test-1, got %s", retrieved.ID)
	}
}

func TestGetByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSqliteDisplayRepository(db)

	_, err := repo.GetByID("nonexistent")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestUpdate(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSqliteDisplayRepository(db)
	display := createTestDisplay("test-1")

	if err := repo.Create(display); err != nil {
		t.Fatalf("failed to create display: %v", err)
	}

	display.Content.Rows[0][0] = "X"

	err := repo.Update(display)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	retrieved, err := repo.GetByID("test-1")
	if err != nil {
		t.Fatalf("failed to retrieve updated display: %v", err)
	}

	if retrieved.Content.Rows[0][0] != "X" {
		t.Errorf("expected first cell to be 'X', got '%s'", retrieved.Content.Rows[0][0])
	}
}

func TestDelete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSqliteDisplayRepository(db)
	display := createTestDisplay("test-1")

	if err := repo.Create(display); err != nil {
		t.Fatalf("failed to create display: %v", err)
	}

	err := repo.Delete("test-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = repo.GetByID("test-1")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestList(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSqliteDisplayRepository(db)

	display1 := createTestDisplay("test-1")
	display2 := createTestDisplay("test-2")

	if err := repo.Create(display1); err != nil {
		t.Fatalf("failed to create display1: %v", err)
	}
	if err := repo.Create(display2); err != nil {
		t.Fatalf("failed to create display2: %v", err)
	}

	displays, err := repo.List()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(displays) != 2 {
		t.Errorf("expected 2 displays, got %d", len(displays))
	}
}
