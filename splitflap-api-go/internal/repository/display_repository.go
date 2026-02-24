package repository

import "splitflap-api-go/internal/model"

// DisplayRepository defines the interface for display data persistence
type DisplayRepository interface {
	// GetByID retrieves a display by its ID
	GetByID(id string) (*model.Display, error)

	// Create persists a new display
	Create(display *model.Display) error

	// Update modifies an existing display
	Update(display *model.Display) error

	// Delete removes a display by its ID
	Delete(id string) error

	// List retrieves all displays
	List() ([]*model.Display, error)
}
