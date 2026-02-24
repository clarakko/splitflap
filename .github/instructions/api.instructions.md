# Go API Coding Standards

## Project Structure

**IMPORTANT:** This project is a monorepo. The Go API is in the `splitflap-api-go/` directory.

### Running Go Commands

**ALWAYS** navigate to the `splitflap-api-go/` directory before running Go commands:

```bash
# ✅ Correct - Always cd to splitflap-api-go first
cd /Users/clara.devers/workspace/catalyst/splitflap/splitflap-api-go
go build ./cmd/api
go test ./...
go run ./cmd/api/main.go

# ❌ Wrong - Running from repo root will fail
go build ./cmd/api  # Will look for wrong path!
```

## General Principles

- Follow Go conventions from effective_go.md
- Prefer composition over inheritance
- Use interfaces for loose coupling
- Keep functions focused and testable
- Write clear, idiomatic code

## Code Style

### Naming Conventions

- Types/Functions: PascalCase for exported (`GetDisplay`), camelCase for unexported (`getDisplay`)
- Constants: PascalCase for exported (`MaxRowCount`)
- Packages: single lowercase word (`handler`, `service`, `model`)
- Interfaces: PascalCase ending in -er (`Fetcher`, `Writer`)

### File Organization

```go
// 1. Package declaration
package handler

// 2. Imports (organized: stdlib, third-party, project)
import (
	"encoding/json"
	"net/http"

	"splitflap/internal/service"
)

// 3. Types
type DisplayHandler struct {
	service *service.DisplayService
}

// 4. Functions
func (h *DisplayHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Implementation
}
```

## Go Patterns

### Handlers

- Implement `http.Handler` interface or use `http.HandlerFunc`
- Accept `http.ResponseWriter` and `*http.Request`
- Extract path variables manually or with string manipulation
- Keep handlers thin - delegate to services

```go
type DisplayHandler struct {
	service *service.DisplayService
}

func (h *DisplayHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/displays/")

	display := h.service.GetDisplay(id)
	if display == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(display)
}
```

### Services

- Contain business logic
- Return domain objects or nil (not error for missing data)
- Use dependency injection via constructor

```go
type DisplayService struct {
	// Dependencies
}

func (s *DisplayService) GetDisplay(id string) *Display {
	// Business logic here
	return nil
}
```

### Models

- Use structs with JSON tags
- Keep in dedicated `model` package
- Match API response structure exactly

```go
type Display struct {
	ID      string          `json:"id"`
	Content DisplayContent  `json:"content"`
	Config  DisplayConfig   `json:"config"`
}

type DisplayContent struct {
	Rows [][]string `json:"rows"`
}

type DisplayConfig struct {
	RowCount int `json:"rowCount"`
	ColCount int `json:"columnCount"`
}
```

## Error Handling

### Phase 1

- Return 404 for missing resources
- Log errors and return 500 for unexpected errors
- Keep error handling simple

````go
if display == nil {
	w.WriteHeader(http.StatusNotFound)
	return
}

if err != nil {
	w.WriteHeader(http.StatusInternalServerError)
	log.Printf("Error: %v", err)
	returtesting package from stdlib
- Test file naming: `{filename}_test.go`
- Use table-driven tests
- Mock dependencies via interfaces

```go
func TestDisplayHandler_GetExisting(t *testing.T) {
	tests := []struct {
		name           string
		displayID      string
		expectedStatus int
	}{
		{"demo exists", "demo", http.StatusOK},
		{"unknown display", "unknown", http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &service.DisplayService{}
			handler := &DisplayHandler{service: service}

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/displays/"+tt.displayID, nil)

			handler.ServeHTTP(w, r)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected %d, got %d", tt.expectedStatus, w.Code)
			}
		})

    @MockBean
    private lateinit var service: DisplayService

    @Test
    fun `GET existing display returns 200`() {
        // Arrange
        val display = Display(...)
        whenever(service.getDisplay("demo")).thenReturn(display)

        // Act & Assert
        mockMvc.get("/api/v1/displays/demo")
            .andExpect { status { isOk() } }
            .andExpect { jsonPath("$.id") { value("demo") } }
    }
}main.go

- Use constants for configuration (Phase 1)
- Pass config via struct fields to handlers (Phase 2)
- Use environment variables for sensitive data (Phase 2+)

```go
const (
	Port = ":8080"
	AllowedOrigin = "http://localhost:5173"
)

func main() {
	handler := &DisplayHandler{
		service: &DisplayService{},
	}

	http.ListenAndServe(Port, corsMiddleware(handler))
}
````

### CORS Middleware

- Create simple middleware for CORS headers
- Keep Phase 1 minimal

```go
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
		Go stdlib only (`net/http`, `encoding/json`, `testing`)
- No external packages

### Don't Add Yet

- ❌ Database drivers (Phase 2)
- ❌ WebSocket libraries (Phase 5)
- ❌ Authentication packages (Phase 8)
- ❌ Third-party frameworks
    }
}
```

## Dependencies

### Phase 1 Only Use

- `spring-boot-starter-web`
- `jackson-module-kotlin`
- `sprcomment on type/function for exported items

````go
// ❌ Bad: Obvious
// Get display by ID
func (s *DisplayService) GetDisplay(id string) *Display

// ✅ Good: Explains decision
// GetDisplay returns the display with the given ID.
// Phase 1: Hardcoded data. Phase 2 will query database.
func (s *DisplayService) GetDisplay(id string) *Display {
	return demoDisplay
Go-Specific

### Nil Safety

- Check for nil explicitly
- Use `if err != nil` pattern for errors
- Return nil for missing values in Phase 1

```go
// ✅ Good
display := s.GetDisplay(id)
if display == nil {
	w.WriteHeader(http.StatusNotFound)
	return
}

w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(display)
````

### Interfaces

- Use interfaces for loose coupling
- Define interfaces where they're used
- Keep interfaces small (1-3 methods)

```go
type DisplayFetcher interface {
	GetDisplay(id string) *Display
}
```

### Error Handling

- Use explicit error returns
- Don't panic in production code
- Log errors before returning
  return ResponseEntity.notFound().build()
  }

```

### Data Classes

- Use for DTOs and value objects
- Automatically get `equals()`, `hashCode()`, `toString()`
- Use `copy()` for immutable updates

### Extension Functions

- Use sparingly
- Only when it truly extends a type's interface
- Don't pollute standard library types

## Git Commits

Use conventional commits:

```

feat(api): add GET /v1/displays/{id} endpoint
fix(api): return 404 for missing displays
test(api): add controller tests for display endpoint
refactor(api): extract display validation logic
docs(api): update API.md with error responses

````

## Phase Discipline

**DO:**
- ✅ Implement only Phase 1 features
- ✅ Add TODO comments referencing future phases
- ✅ Keep code simple and readable

**DON'T:**
- ❌ Add database code (Phase 2)
- ❌ Add POST/PUT/DELETE endpoints (Phase 2)
- ❌ Add authentication (Phase 8)
- ❌ Over-engineer for future needs

```kotlin
// ✅ Good: Simple Phase 1 implementation
private val demoDisplay = Display(
    id = "demo",
    content = DisplayContent(rows = listOf(...))
)

// ❌ Bad: Premature abstraction for Phase 2
interface DisplayRepository {
    fun findById(id: String): Display?
}
````
