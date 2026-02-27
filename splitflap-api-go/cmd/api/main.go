package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"splitflap-api-go/internal/handler"
	"splitflap-api-go/internal/middleware"
	"splitflap-api-go/internal/model"
	"splitflap-api-go/internal/repository"
	"splitflap-api-go/internal/service"

	_ "github.com/mattn/go-sqlite3"
)

const schemaSQL = `-- SplitFlap Display Database Schema
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

func main() {
	// Initialize database
	db, err := initDatabase()
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()

	// Create repository
	displayRepo := repository.NewSqliteDisplayRepository(db)

	// Seed demo display if database is empty
	if err := seedDemoDisplay(displayRepo); err != nil {
		log.Fatalf("failed to seed demo display: %v", err)
	}

	// Create service and handler
	displayService := service.NewDisplayService(displayRepo)
	displayHandler := handler.NewDisplayHandler(displayService)

	mux := http.NewServeMux()
	// Register both patterns to handle requests with and without trailing slash
	mux.Handle("/api/v1/displays", displayHandler)
	mux.Handle("/api/v1/displays/", displayHandler)

	addr := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}

	server := &http.Server{
		Addr:    addr,
		Handler: middleware.WithCORS(mux),
	}

	log.Printf("Splitflap API listening on %s", addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server failed: %v", err)
	}
}

func initDatabase() (*sql.DB, error) {
	// Create data directory if it doesn't exist
	dataDir := "./data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}

	dbPath := filepath.Join(dataDir, "splitflap.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Run schema
	if _, err := db.Exec(schemaSQL); err != nil {
		db.Close()
		return nil, err
	}

	log.Printf("Database initialized at %s", dbPath)
	return db, nil
}

func seedDemoDisplay(repo repository.DisplayRepository) error {
	// Check if demo display already exists
	existing, err := repo.GetByID("demo")
	if err == nil && existing != nil {
		log.Println("Demo display already exists, skipping seed")
		return nil
	}
	if err != nil && err != repository.ErrNotFound {
		return err
	}

	// Create demo display
	demoDisplay := &model.Display{
		ID: "demo",
		Content: model.DisplayContent{
			Rows: [][]string{
				{"H", "E", "L", "L", "O", " ", "W", "O", "R", "L"},
				{"D", " ", "S", "P", "L", "I", "T", "F", "L", "A"},
				{"P", " ", "D", "I", "S", "P", "L", "A", "Y", " "},
				{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"},
				{"-", ":", ".", ",", " ", "A", "Z", "a", "z", "!"},
			},
		},
		Config: model.DisplayConfig{
			RowCount:    5,
			ColumnCount: 10,
		},
	}

	if err := repo.Create(demoDisplay); err != nil {
		return err
	}

	log.Println("Demo display seeded successfully")
	return nil
}
