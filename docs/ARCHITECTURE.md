# Architecture

## Overview

SplitFlap is a three-tier web application for creating and embedding animated split-flap displays.

```
┌─────────────────┐
│  splitflap-web  │  React builder app (Phase 2+)
│   (Port 5173)   │
└────────┬────────┘
         │ HTTP REST
         ▼
┌─────────────────┐
│ splitflap-api   │  Kotlin/Spring Boot backend
│   (Port 8080)   │
└────────┬────────┘
         │ (Phase 2: JPA)
         ▼
┌─────────────────┐
│   Database      │  H2 (dev) / Postgres (prod)
│                 │  Not used in Phase 1
└─────────────────┘

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
- ✅ Kotlin/Spring Boot REST controller
- ✅ JSON response with display grid data
- ✅ Basic error handling (404 for missing display)
- ✅ CORS configuration for local dev

### Out of Scope

- ❌ Database persistence (Phase 2)
- ❌ Display CRUD operations (Phase 2)
- ❌ React builder UI (Phase 2-3)
- ❌ Web component embed (Phase 4)
- ❌ Real-time updates (Phase 5)
- ❌ External data sources (Phase 6)
- ❌ Authentication (Phase 8)

---

## Technology Stack

### Backend (splitflap-api)

| Technology | Version | Purpose |
|------------|---------|---------|
| Kotlin | 2.2.21 | Primary language |
| Spring Boot | 4.0.0 | Web framework |
| Spring Web MVC | (included) | REST endpoints |
| Jackson Kotlin | (included) | JSON serialization |
| Spring Data JPA | (included) | Phase 2: Database access |
| H2 Database | (included) | Phase 2: In-memory dev database |
| JUnit 5 | (included) | Testing |

### Frontend (splitflap-web) - Phase 2+

| Technology | Version | Purpose |
|------------|---------|---------|
| React | TBD | UI framework |
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
dev.clarakko.splitflap_api/
├── SplitflapApiApplication.kt       # Spring Boot entry point
├── controller/
│   └── DisplayController.kt         # REST endpoints
├── service/
│   └── DisplayService.kt            # Business logic (in-memory data)
├── dto/
│   ├── Display.kt                   # Response DTO
│   ├── DisplayContent.kt            # Content DTO
│   └── DisplayConfig.kt             # Config DTO
└── config/
    └── CorsConfig.kt                # CORS settings
```

### Layer Responsibilities

#### Controller Layer
- HTTP request/response handling
- Path variable extraction (`/displays/{id}`)
- Status code mapping (200, 404, 500)
- Exception handling

#### Service Layer
- Business logic (currently minimal)
- Data retrieval (hardcoded in Phase 1)
- Phase 2: Database interaction via repositories

#### DTO Layer
- Data transfer objects matching API spec
- JSON serialization via Jackson
- No business logic

### Data Flow (Phase 1)

```
1. GET /api/v1/displays/demo
   ↓
2. DisplayController.getDisplay("demo")
   ↓
3. DisplayService.getDisplay("demo")
   ↓
4. Return hardcoded Display object
   ↓
5. Jackson serializes to JSON
   ↓
6. HTTP 200 with JSON body
```

### Error Handling

```kotlin
@RestController
@RequestMapping("/api/v1/displays")
class DisplayController(private val service: DisplayService) {
    
    @GetMapping("/{id}")
    fun getDisplay(@PathVariable id: String): ResponseEntity<Display> {
        val display = service.getDisplay(id)
        return if (display != null) {
            ResponseEntity.ok(display)
        } else {
            ResponseEntity.notFound().build()
        }
    }
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

### application.yaml (Phase 1)

```yaml
spring:
  application:
    name: splitflap-api
  
  # Phase 2: Database config will go here
  # datasource:
  #   url: jdbc:h2:mem:splitflap
  #   driver-class-name: org.h2.Driver
  
server:
  port: 8080

# Phase 1: Enable CORS for local frontend dev
# (Will be configured via CorsConfig.kt)
```

### CORS Configuration (Phase 1)

```kotlin
@Configuration
class CorsConfig : WebMvcConfigurer {
    override fun addCorsMappings(registry: CorsRegistry) {
        registry.addMapping("/api/v1/**")
            .allowedOrigins("http://localhost:5173", "http://localhost:3000")
            .allowedMethods("GET")
            .allowedHeaders("*")
    }
}
```

---

## Evolution Through Phases

### Phase 1 → Phase 2
**Changes:**
- Add `@Entity` to Display model
- Add `@Repository` interface for DisplayRepository
- Add H2 database configuration
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

```kotlin
@WebMvcTest(DisplayController::class)
class DisplayControllerTest {
    
    @Test
    fun `GET demo display returns 200`()
    
    @Test
    fun `GET unknown display returns 404`()
    
    @Test
    fun `Response matches JSON schema`()
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

## Deployment (Future)

### Phase 1-2
- Local development only
- `./gradlew bootRun` for backend
- `npm run dev` for frontend

### Phase 8+
- Backend: Docker container → Kubernetes/Railway/Fly.io
- Frontend: Static build → Vercel/Netlify/Cloudflare Pages
- Database: Managed Postgres (Supabase/Neon/RDS)
- Embed script: CDN (Cloudflare/Fastly)
