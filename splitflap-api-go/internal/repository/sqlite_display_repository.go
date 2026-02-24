package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"splitflap-api-go/internal/model"

	_ "github.com/mattn/go-sqlite3"
)

var (
	ErrNotFound      = errors.New("display not found")
	ErrAlreadyExists = errors.New("display already exists")
)

// SqliteDisplayRepository implements DisplayRepository using SQLite
type SqliteDisplayRepository struct {
	db *sql.DB
}

// NewSqliteDisplayRepository creates a new SQLite-backed repository
func NewSqliteDisplayRepository(db *sql.DB) *SqliteDisplayRepository {
	return &SqliteDisplayRepository{db: db}
}

// GetByID retrieves a display by its ID
func (r *SqliteDisplayRepository) GetByID(id string) (*model.Display, error) {
	query := `
		SELECT id, content_rows, row_count, column_count
		FROM displays
		WHERE id = ?
	`

	var display model.Display
	var contentRowsJSON string

	err := r.db.QueryRow(query, id).Scan(
		&display.ID,
		&contentRowsJSON,
		&display.Config.RowCount,
		&display.Config.ColumnCount,
	)

	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query display: %w", err)
	}

	// Deserialize content rows from JSON
	if err := json.Unmarshal([]byte(contentRowsJSON), &display.Content.Rows); err != nil {
		return nil, fmt.Errorf("failed to unmarshal content rows: %w", err)
	}

	return &display, nil
}

// Create persists a new display
func (r *SqliteDisplayRepository) Create(display *model.Display) error {
	// Serialize content rows to JSON
	contentRowsJSON, err := json.Marshal(display.Content.Rows)
	if err != nil {
		return fmt.Errorf("failed to marshal content rows: %w", err)
	}

	query := `
		INSERT INTO displays (id, content_rows, row_count, column_count, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	_, err = r.db.Exec(query,
		display.ID,
		string(contentRowsJSON),
		display.Config.RowCount,
		display.Config.ColumnCount,
		now,
		now,
	)

	if err != nil {
		// SQLite constraint error for duplicate ID
		if err.Error() == "UNIQUE constraint failed: displays.id" {
			return ErrAlreadyExists
		}
		return fmt.Errorf("failed to insert display: %w", err)
	}

	return nil
}

// Update modifies an existing display
func (r *SqliteDisplayRepository) Update(display *model.Display) error {
	// Serialize content rows to JSON
	contentRowsJSON, err := json.Marshal(display.Content.Rows)
	if err != nil {
		return fmt.Errorf("failed to marshal content rows: %w", err)
	}

	query := `
		UPDATE displays
		SET content_rows = ?, row_count = ?, column_count = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.Exec(query,
		string(contentRowsJSON),
		display.Config.RowCount,
		display.Config.ColumnCount,
		time.Now(),
		display.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update display: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// Delete removes a display by its ID
func (r *SqliteDisplayRepository) Delete(id string) error {
	query := `DELETE FROM displays WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete display: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// List retrieves all displays
func (r *SqliteDisplayRepository) List() ([]*model.Display, error) {
	query := `
		SELECT id, content_rows, row_count, column_count
		FROM displays
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query displays: %w", err)
	}
	defer rows.Close()

	var displays []*model.Display

	for rows.Next() {
		var display model.Display
		var contentRowsJSON string

		err := rows.Scan(
			&display.ID,
			&contentRowsJSON,
			&display.Config.RowCount,
			&display.Config.ColumnCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan display row: %w", err)
		}

		// Deserialize content rows from JSON
		if err := json.Unmarshal([]byte(contentRowsJSON), &display.Content.Rows); err != nil {
			return nil, fmt.Errorf("failed to unmarshal content rows: %w", err)
		}

		displays = append(displays, &display)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating display rows: %w", err)
	}

	return displays, nil
}
