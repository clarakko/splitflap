# Architecture

## Overview

SplitFlap is a three-tier web application for creating and embedding animated split-flap displays.

```
┌───────────────────────┐
│ splitflap-web-solid   │  SolidJS builder app (Phase 2+)
│    (Port 5173)        │
└──────────┬────────────┘
           │ HTTP REST
           ▼
┌───────────────────────┐
│  splitflap-api-go     │  Go stdlib backend
│    (Port 8080)        │
└──────────┬────────────┘
           │ SQL (Phase 2+)
           ▼
┌───────────────────────┐
│   Database            │  SQLite (data/splitflap.db)
│                       │  Repository pattern enables Postgres migration
└───────────────────────┘

┌─────────────────┐
│splitflap-embed  │  Web component (Phase 4)
│ (Vanilla JS)    │  Fetches from API, renders display
└─────────────────┘
```

---

## Phase 1: MVP Architecture

### Goals

1. Prove split-flap flip animation works in browser
2. Establish API contract for display data
3. Validate data model for grid-based displays
4. Set foundation for future phases

### In Scope

- ✅ Single hardcoded display served via REST API
- ✅ Go REST controller
- ✅ JSON response with display grid data
- ✅ Basic error handling (404 for missing display)
- ✅ CORS configuration for local dev

### Out of Scope

- ✅ Database persistence (Phase 2 - SQLite with repository pattern)
- ❌ Display CRUD operations (Phase 2)
- ❌ SolidJS builder UI (Phase 2-3)
- ❌ Web component embed (Phase 4)
- ❌ Real-time updates (Phase 5)
- ❌ External data sources (Phase 6)
- ❌ Authentication (Phase 8)

---

## Technology Stack

### Backend (splitflap-api-go)

| Technology | Version | Purpose |
|------------|---------|---------|
| Go | 1.22 | Primary language |
| net/http | (stdlib) | REST endpoints |
| encoding/json | (stdlib) | JSON serialization |
| testing | (stdlib) | Unit and contract tests |

### Frontend (splitflap-web) - Phase 2+

| Technology | Version | Purpose |
|------------|---------|---------|
| SolidJS | TBD | UI framework |
| Vite | TBD | Build tool |
| TypeScript | TBD | Type safety |

### Embed (splitflap-embed) - Phase 4

| Technology | Purpose |
|------------|---------|
| Vanilla JavaScript | Minimal dependencies |
| Web Components | Standards-based embedding |

---

## Backend Architecture (Phase 1)

### Package Structure

```
splitflap-api-go/
├── cmd/
│   └── api/
│       └── main.go                  # HTTP server entry point
├── internal/
│   ├── handler/
│   │   └── display_handler.go       # REST endpoints
│   ├── service/
│   │   └── display_service.go       # Business logic (in-memory data)
│   ├── model/
│   │   └── display.go               # Response DTO
│   └── middleware/
│       └── cors.go                  # CORS settings
```

### Layer Responsibilities

#### Handler Layer
- HTTP request/response handling
- Path variable extraction (`/displays/{id}`)
- Status code mapping (200, 404, 500)
- Exception handling

#### Service Layer
- Business logic (currently minimal)
- Data retrieval (hardcoded in Phase 1)
- Phase 2: Database interaction via repositories (✅ implemented with SQLite)

#### Model Layer
- Data transfer objects matching API spec
- JSON serialization via encoding/json
- No business logic

### Data Flow (Phase 1)

```
1. GET /api/v1/displays/demo
   ↓
2. Handler routes request to service
   ↓
3. DisplayService.GetDisplay("demo")
   ↓
4. Return hardcoded Display object
   ↓
5. encoding/json serializes to JSON
   ↓
6. HTTP 200 with JSON body
```

### Error Handling

```go
func (h *DisplayHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  id := strings.TrimPrefix(r.URL.Path, "/api/v1/displays/")
  display := h.service.GetDisplay(id)
  if display == nil {
    w.WriteHeader(http.StatusNotFound)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  json.NewEncoder(w).Encode(display)
}
```

---

## Frontend Architecture (Phase 2+)

### Component Hierarchy (Planned)

```
<App>
  ├── <DisplayPreview>              # Phase 2: Fetch and render display
  │     └── <SplitFlapBoard>        # Phase 2: Grid container
  │           └── <SplitFlapCell>   # Phase 1-2: Single flap animation
  │                 └── <Flap>      # Phase 2: Individual character
  │
  ├── <DisplayBuilder>              # Phase 3: Configuration UI
  │     ├── <GridEditor>
  │     └── <ConfigPanel>
  │
  └── <EmbedCodeGenerator>          # Phase 4: Generate embed snippet
```

### Key Frontend Concepts

#### Mechanical Simulation

Split-flap displays are **mechanical devices** with physical constraints:

1. **Sequential flipping**: Characters flip through the full sequence (A→B→C→...→Z→0→1→...→9→ space)
2. **No instant changes**: Cannot jump from 'A' to 'Z' without flipping through B, C, D, etc.
3. **Timing**: Each flip takes ~50-100ms (configurable)
4. **Synchronization**: Multiple cells can flip simultaneously but follow same rules
5. **Character set**: Physical devices have fixed character drums (Phase 1: alphanumeric only)

**Phase 3+ Character Set Extensions:**

- Image-based sets (airline logos, weather icons, custom images)
- Requires sprite sheet loading and CSS positioning instead of text rendering
- Cell animation still sequential: flips through image sequence
- Example: Airline logos flip through: `[AA → BA → DL → UA → ...]`

**Implementation implications:**
- Frontend calculates shortest path through character set
- CSS animations simulate rotation timing
- State management tracks current vs. target character per cell
- Animation queue handles rapid updates

**Example**: Changing "A" to "D" requires 3 flips: A→B→C→D

#### Animation State (Frontend Responsibility)

```typescript
interface CellState {
  current: string;      // Currently displayed character
  target: string;       // Target character from API
  isFlipping: boolean;  // Animation in progress
  progress: number;     // 0.0 to 1.0
}
```

Backend serves **target state** only. Frontend owns **animation state**.

---

## Configuration

### CORS Configuration (Phase 1)

CORS is configured in the Go handler middleware to allow requests from the local Vite dev server:

```go
w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
```

---

## Evolution Through Phases

### Phase 1 → Phase 2
**Changes:**
- Add database model tags to Display struct
- Create DisplayRepository interface
- Add Postgres connection pool configuration
- Seed database with demo display
- Service layer queries database instead of returning hardcoded data

**Backward compatibility:** API contract unchanged

---

### Phase 2 → Phase 3

**Changes:**

- Add POST, PUT, DELETE endpoints
- Add builder UI in splitflap-web
- Extend DisplayConfig with `flipSpeed`, `characterSet` (enum: alphanumeric, numeric, custom)
- Add optional `metadata` to DisplayContent
- Support custom character set definitions (text or image-based)
- Add character set asset management for image-based sets

**Backward compatibility:** New fields are optional; existing GET endpoint works; default characterSet is "alphanumeric"

---

### Phase 3 → Phase 4
**Changes:**
- Create splitflap-embed web component
- Component calls same GET endpoint
- Add embed code generator in builder UI

**Backward compatibility:** No API changes

---

### Phase 4 → Phase 5
**Changes:**
- Add WebSocket endpoint for real-time updates
- Modify frontend to subscribe to updates
- Add PATCH endpoint for partial updates

**Backward compatibility:** GET endpoint still works for polling clients

---

### Phase 5 → Phase 6
**Changes:**
- Add DataSource configuration to Display entity
- Add background job to fetch external API data
- Update display content periodically
- GET endpoint returns latest fetched data

**Backward compatibility:** API response format unchanged

---

## Design Principles

### 1. Spec-Driven Development
- Document data model and API before implementing
- Frontend consumes API contract, not backend implementation details
- Changes require spec update first

### 2. Phase Discipline
- Only implement features in current phase
- Resist scope creep ("wouldn't it be cool if...")
- Each phase delivers working, testable functionality

### 3. Separation of Concerns
- Backend owns data, frontend owns presentation
- Backend serves target state, frontend manages animation
- Embed component is independent of builder UI

### 4. Progressive Enhancement
- Phase 1 proves core concept (flip animation + API)
- Each phase adds one major capability
- No "big bang" deployment

### 5. Backward Compatibility
- API versioned at path level (`/api/v1/`, `/api/v2/`, etc.)
- Breaking changes require new version; existing versions maintained
- Optional fields for new features within same version
- Database migrations (Phase 2+)
- Version negotiation via path (simple, explicit, cacheable)

---

## Security Considerations

### Phase 1
- None (local development only)
- CORS allows localhost origins

### Phase 8 (Future)
- JWT-based authentication
- User-scoped displays
- Public vs. private displays
- Rate limiting
- Input validation/sanitization

---

## Performance Considerations

### Phase 1
- Not applicable (single hardcoded display)

### Phase 2+
- Database indexing on display ID
- JSON response caching headers
- Phase 5: WebSocket reduces polling overhead
- Phase 6: Cache external API responses

---

## Testing Strategy

### Backend Testing (Phase 1)

```go
func TestGetDisplay(t *testing.T) {
    tests := []struct {
        name           string
        displayID      string
        expectedStatus int
    }{
        {"get demo display", "demo", http.StatusOK},
        {"get unknown display", "unknown", http.StatusNotFound},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Frontend Testing (Phase 2+)
- Component unit tests (Jest/Vitest)
- Animation state machine tests
- API integration tests (mock responses)
- E2E tests (Playwright/Cypress)

---

## Development Workflow

### Phase 1

1. ✅ Document data model, API, architecture
2. ⏳ Implement backend controller + service
3. ⏳ Add CORS configuration
4. ⏳ Test GET endpoint manually (curl)
5. ⏳ Write unit tests
6. ⏳ Frontend: Build flip animation component
7. ⏳ Frontend: Fetch from API and render

### Phase 2+

Follow spec-first approach:
1. Update docs with new features
2. Update STATUS.md with current task
3. Implement backend changes
4. Implement frontend changes
5. Test integration
6. Update ROADMAP checkboxes

---

## go run ./cmd/api/main.go` for backend
- `pnpm
### Phase 1-2
- Local development only
- `./gradlew bootRun` for backend
- `npm run dev` for frontend

### Phase 8+
- Backend: Docker container → Kubernetes/Railway/Fly.io
- Frontend: Static build → Vercel/Netlify/Cloudflare Pages
- Database: Managed Postgres (Supabase/Neon/RDS)
- Embed script: CDN (Cloudflare/Fastly)
