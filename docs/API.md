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

### GET /api/v1/displays/{id}

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
      ["10:45", "NEW YORK", "5", "DELAYED"],
      ["11:00", "PHILADELPHIA", "7", "BOARDING"],
      ["11:15", "WASHINGTON", "2", "ON TIME"]
    ]
  },
  "config": {
    "rowCount": 5,
    "columnCount": 4
  }
}
```

**Status: 404 Not Found**

```
(no response body)
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

## Phase 2: Persistence & CRUD

Database-backed persistence with full CRUD operations.

### GET /api/v1/displays

List all displays.

#### Response

**Status: 200 OK**

```json
[
  {
    "id": "demo",
    "content": {
      "rows": [["H","E","L","L","O"," ","W","O","R","L"],["D"," ","S","P","L","I","T","F","L","A"],["P"," ","D","I","S","P","L","A","Y"," "],["0","1","2","3","4","5","6","7","8","9"],["-",":",".","," ","A","Z","a","z","!"]]
    },
    "config": {
      "rowCount": 5,
      "columnCount": 10
    }
  }
]
```

#### Example Request

```bash
curl http://localhost:8080/api/v1/displays/
```

---

### POST /api/v1/displays

Create a new display.

#### Request Body

```json
{
  "id": "my-display",
  "content": {
    "rows": [["H","E","L","L","O"]]
  },
  "config": {
    "rowCount": 1,
    "columnCount": 5
  }
}
```

#### Validation Rules

- `id`: Required, must be unique
- `rowCount`: 1-20 (inclusive)
- `columnCount`: 1-10 (inclusive)
- `content.rows` length must match `rowCount`
- Each row length must match `columnCount`
- Each cell must be a single character

#### Response

**Status: 201 Created**

```json
{
  "id": "my-display",
  "content": {
    "rows": [["H","E","L","L","O"]]
  },
  "config": {
    "rowCount": 1,
    "columnCount": 5
  }
}
```

**Status: 400 Bad Request**

```json
{
  "error": "rowCount must be between 1 and 20"
}
```

**Status: 409 Conflict**

```json
{
  "error": "display already exists"
}
```

#### Example Request

```bash
curl -X POST http://localhost:8080/api/v1/displays/ \
  -H "Content-Type: application/json" \
  -d '{"id":"test","content":{"rows":[["H","I"]]},"config":{"rowCount":1,"columnCount":2}}'
```

---

### PUT /api/v1/displays/{id}

Update an existing display.

#### Parameters

| Parameter | Type | Location | Required | Description |
|-----------|------|----------|----------|-------------|
| `id` | String | Path | Yes | Display identifier (must match body) |

#### Request Body

Same as POST. The `id` in the body must match the `id` in the URL path.

#### Response

**Status: 200 OK**

Returns the updated display.

**Status: 400 Bad Request**

```json
{
  "error": "ID mismatch"
}
```

**Status: 404 Not Found**

Display does not exist.

#### Example Request

```bash
curl -X PUT http://localhost:8080/api/v1/displays/test \
  -H "Content-Type: application/json" \
  -d '{"id":"test","content":{"rows":[["B","Y","E"]]},"config":{"rowCount":1,"columnCount":3}}'
```

---

### DELETE /api/v1/displays/{id}

Delete a display.

#### Parameters

| Parameter | Type | Location | Required | Description |
|-----------|------|----------|----------|-------------|
| `id` | String | Path | Yes | Display identifier |

#### Response

**Status: 204 No Content**

Display successfully deleted (no response body).

**Status: 404 Not Found**

Display does not exist.

#### Example Request

```bash
curl -X DELETE http://localhost:8080/api/v1/displays/test
```

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

All errors follow this format in Phase 2+:

```json
{
  "error": "Human-readable error message",
  "details": "Optional additional context"
}
```

### HTTP Status Codes

| Code | Meaning | When Used |
|------|---------|-----------|
| 200 | OK | Successful GET or PUT |
| 201 | Created | Successful POST (resource created) |
| 204 | No Content | Successful DELETE |
| 400 | Bad Request | Validation error or malformed request |
| 404 | Not Found | Display ID doesn't exist |
| 409 | Conflict | Attempt to create display with duplicate ID |
| 500 | Internal Server Error | Unexpected server error |

---

## Future Endpoints (Phase 3+)

Not yet implemented:

- `PATCH /api/v1/displays/{id}` - Partial update (Phase 5)
- WebSocket `/ws/api/v1/displays/{id}` - Real-time updates (Phase 5)

---

## CORS Configuration

Phase 1-2 enables CORS for local development:

- Allow origins: `http://localhost:3000`, `http://localhost:5173` (Vite default)
- Allow methods: `GET`, `POST`, `PUT`, `DELETE`
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

```go
// Handler
mux.Handle("/api/v1/displays/", displayHandler)

// Service (in-memory data)
type DisplayService struct {
  demoDisplay model.Display
}

func (s *DisplayService) GetDisplay(id string) *model.Display

// DTOs
type Display struct {
  ID      string         `json:"id"`
  Content DisplayContent `json:"content"`
  Config  DisplayConfig  `json:"config"`
}

type DisplayContent struct {
  Rows [][]string `json:"rows"`
}

type DisplayConfig struct {
  RowCount    int `json:"rowCount"`
  ColumnCount int `json:"columnCount"`
}
```

### Hardcoded Demo Data

```go
demoDisplay := model.Display{
  ID: "demo",
  Content: model.DisplayContent{
    Rows: [][]string{
      {"TIME", "DESTINATION", "PLATFORM", "STATUS"},
      {"10:30", "BOSTON", "3", "ON TIME"},
      {"10:45", "NEW YORK", "5", "DELAYED"},
      {"11:00", "PHILADELPHIA", "7", "BOARDING"},
      {"11:15", "WASHINGTON", "2", "ON TIME"},
    },
  },
  Config: model.DisplayConfig{RowCount: 5, ColumnCount: 4},
}
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
