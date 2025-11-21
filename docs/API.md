# API Specification

## Phase 1: MVP

Single REST endpoint serving hardcoded display data.

### Base URL

```
http://localhost:8080/api/v1
```

**API Version**: v1

All endpoints are versioned to allow backward-compatible evolution. Future breaking changes will increment the version (v2, v3, etc.).

### Authentication

None (Phase 8 will add authentication)

---

## Endpoints

### GET /v1/displays/{id}

Retrieve a display configuration and content by ID.

#### Parameters

| Parameter | Type | Location | Required | Description |
|-----------|------|----------|----------|-------------|
| `id` | String | Path | Yes | Display identifier |

#### Response

**Status: 200 OK**

```json
{
  "id": "demo",
  "content": {
    "rows": [
      ["TIME", "DESTINATION", "PLATFORM", "STATUS"],
      ["10:30", "BOSTON", "3", "ON TIME"],
      ["10:45", "NEW YORK", "5", "DELAYED"]
    ]
  },
  "config": {
    "rowCount": 3,
    "columnCount": 4
  }
}
```

**Status: 404 Not Found**

```json
{
  "error": "Display not found",
  "id": "unknown-id"
}
```

#### Example Request

```bash
curl http://localhost:8080/api/v1/displays/demo
```

#### Notes

- **Phase 1**: Only one display exists with ID `"demo"`
- **Phase 2**: Will support multiple displays stored in database
- Content is static (hardcoded in backend)
- No caching headers yet (Phase 5 will add real-time updates)

---

## Data Types

### Display

| Field | Type | Description |
|-------|------|-------------|
| `id` | String | Unique display identifier |
| `content` | DisplayContent | Display content grid |
| `config` | DisplayConfig | Display configuration |

### DisplayContent

| Field | Type | Description |
|-------|------|-------------|
| `rows` | String[][] | 2D array of cell values |

### DisplayConfig

| Field | Type | Description |
|-------|------|-------------|
| `rowCount` | Integer | Number of rows |
| `columnCount` | Integer | Number of columns |

---

## Error Responses

All errors follow this format:

```json
{
  "error": "Human-readable error message",
  "details": "Optional additional context"
}
```

### HTTP Status Codes

| Code | Meaning | When Used |
|------|---------|-----------|
| 200 | OK | Successful GET |
| 404 | Not Found | Display ID doesn't exist |
| 500 | Internal Server Error | Unexpected server error |

---

## Future Endpoints (Phase 2+)

Not implemented in Phase 1:

- `POST /v1/displays` - Create new display
- `PUT /v1/displays/{id}` - Update display content
- `DELETE /v1/displays/{id}` - Delete display
- `GET /v1/displays` - List all displays
- `PATCH /v1/displays/{id}` - Partial update (Phase 5)
- WebSocket `/ws/v1/displays/{id}` - Real-time updates (Phase 5)

---

## CORS Configuration

Phase 1 will enable CORS for local development:

- Allow origins: `http://localhost:3000`, `http://localhost:5173` (Vite default)
- Allow methods: `GET`
- Allow headers: `Content-Type`

Phase 3 will add configurable CORS for production deployments.

---

## Content-Type

All requests and responses use:

```
Content-Type: application/json
```

---

## Implementation Notes

### Phase 1 Backend Structure

```kotlin
// Controller
@RestController
@RequestMapping("/api/v1/displays")
class DisplayController(private val displayService: DisplayService)

// Service (in-memory data)
@Service
class DisplayService {
    fun getDisplay(id: String): Display?
}

// DTOs
data class Display(
    val id: String,
    val content: DisplayContent,
    val config: DisplayConfig
)

data class DisplayContent(val rows: List<List<String>>)
data class DisplayConfig(val rowCount: Int, val columnCount: Int)
```

### Hardcoded Demo Data

```kotlin
private val demoDisplay = Display(
    id = "demo",
    content = DisplayContent(
        rows = listOf(
            listOf("TIME", "DESTINATION", "PLATFORM", "STATUS"),
            listOf("10:30", "BOSTON", "3", "ON TIME"),
            listOf("10:45", "NEW YORK", "5", "DELAYED"),
            listOf("11:00", "PHILADELPHIA", "7", "BOARDING"),
            listOf("11:15", "WASHINGTON", "2", "ON TIME")
        )
    ),
    config = DisplayConfig(rowCount = 5, columnCount = 4)
)
```

---

## Testing

### Manual Testing

```bash
# Should return demo display
curl http://localhost:8080/api/v1/displays/demo

# Should return 404
curl http://localhost:8080/api/v1/displays/nonexistent
```

### Unit Test Coverage

- ✅ Controller returns 200 for existing display
- ✅ Controller returns 404 for missing display
- ✅ Response JSON matches schema
- ✅ Service returns demo data correctly
