# Data Model

## Phase 1: MVP (In-Memory)

Phase 1 uses hardcoded, in-memory data. No database persistence yet (see Phase 2).

### Display Entity

Represents a single split-flap display configuration and content.

#### Schema

| Field | Type | Description |
|-------|------|-------------|
| `id` | String | Unique identifier for the display |
| `content` | DisplayContent | The content to display |
| `config` | DisplayConfig | Display configuration settings |

#### DisplayContent

Contains the actual data shown on the display as a 2D grid.

| Field | Type | Description |
|-------|------|-------------|
| `rows` | List<List<String>> | 2D array of strings, each cell is one character or short string |

**Rules (Phase 1):**

- Each string represents one "flap unit" (typically 1-10 characters)
- Alphanumeric characters only: `A-Z`, `0-9`, space, and basic punctuation (`.`, `,`, `:`, `-`)
- Empty cells represented as empty string `""`
- All rows must have same length (columnCount)

**Character Set Evolution (Phase 3+):**

- Phase 3 will add `characterSet` config field (e.g., `"alphanumeric"`, `"numeric"`, `"custom"`)
- Future phases may support image-based character sets (e.g., airline logos, weather icons)
- Image sets require different rendering: sprite sheets or image URLs instead of text
- Cell values would reference image identifiers (e.g., `"LOGO_AA"`, `"ICON_SUNNY"`)

#### DisplayConfig

Configuration for display dimensions and behavior.

| Field | Type | Description |
|-------|------|-------------|
| `rowCount` | Int | Number of rows in the display |
| `columnCount` | Int | Number of columns in the display |

**Phase 1 constraints:**
- Fixed dimensions (no resizing)
- rowCount: 1-20
- columnCount: 1-10

### Example: Departure Board Display

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

### Example: Scoreboard Display

```json
{
  "id": "demo",
  "content": {
    "rows": [
      ["", "Q1", "Q2", "Q3", "Q4", "TOTAL"],
      ["HOME", "21", "14", "", "", "35"],
      ["AWAY", "17", "10", "", "", "27"]
    ]
  },
  "config": {
    "rowCount": 3,
    "columnCount": 6
  }
}
```

## Phase 2+: Evolution Path

Future phases will add:

- **Phase 2**: Database persistence (H2/Postgres), multiple displays, timestamps
- **Phase 3**: Metadata (header rows, column labels, cell styling), configurable character sets (alphanumeric, numeric, custom text, image-based), flip speed settings
- **Phase 4**: No model changes (embed component consumes existing API)
- **Phase 5**: No model changes (real-time updates via WebSocket, same data structure)
- **Phase 6**: DataSource field (external API integration, refresh intervals, field mappings)
- **Phase 7**: Schedule field (time-based content rules)
- **Phase 8**: User ownership (userId, isPublic, tags, createdAt)

### Versioning Strategy

All displays will include a `version` field in future phases:

```json
{
  "version": "1.0",
  "id": "demo",
  "content": {...}
}
```

This allows backward-compatible schema evolution and graceful migrations.

## Design Decisions

### Why 2D Array of Strings?

- **Simplicity**: Matches mental model (grid of cells)
- **Flexibility**: No rigid schema per display "type"
- **Animation-friendly**: Frontend iterates cells, flips each character
- **Resize-friendly**: Adding/removing rows or columns is array manipulation
- **JSON-native**: No complex serialization needed

### Why No Character Validation in Phase 1?

- Reduces API complexity
- Allows experimentation with display content
- Validation can be added in Phase 3 as config option
- Invalid characters can render as blank/fallback on frontend

### Why No Cell Metadata Yet?

- Phase 1 focuses on proving animation works
- Styling, types, and cell-level config add complexity
- Can be added non-breaking in Phase 3 via optional fields

### Why Empty String for Blank Cells?

- Consistent with string type (no nulls to handle)
- Renders naturally as blank flap
- Simpler frontend logic than null checking
