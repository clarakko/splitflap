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
| `rows` | List<List<String>> | 2D array of single-character strings |

**Rules (Phase 1 MVP):**

- **Each cell contains exactly one character**
- Supported characters: `A-Z`, `0-9`, space, and basic punctuation (`.`, `,`, `:`, `-`)
- Empty cells represented as space `" "`
- All rows must have same length (columnCount)
- Total cells = rowCount × columnCount

**Future Enhancement (Phase 3+):**

- Phase 3 will add `characterSet` config field for different character sets
- Phase 3+ may support multi-character cells (words, numbers, symbols on single flap)
- Future phases may support image-based cells (e.g., airline logos, weather icons)
- Image cells would reference identifiers instead of character strings

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

### Example: Phase 1 MVP Display (Single-Character Cells)

```json
{
  "id": "demo",
  "content": {
    "rows": [
      ["H", "E", "L", "L", "O", " ", "W", "O", "R", "L"],
      ["D", " ", "S", "P", "L", "I", "T", "F", "L", "A"],
      ["P", " ", "D", "I", "S", "P", "L", "A", "Y", " "],
      ["0", "1", "2", "3", "4", "5", "6", "7", "8", "9"],
      ["-", ":", ".", ",", " ", "A", "Z", "a", "z", "!"]
    ]
  },
  "config": {
    "rowCount": 5,
    "columnCount": 10
  }
}
```

**Note:** Each cell contains exactly one character. This is the Phase 1 MVP constraint. Future phases may support multi-character cells (Phase 3+) or image-based cells (Phase 4+).

### Future: Departure Board Example (Phase 3+)

Once we support multi-character cells and custom formatting, displays could look like:

```json
{
  "id": "board-1",
  "content": {
    "rows": [
      ["10:30", "BOSTON", "3"],
      ["10:45", "NEW YORK", "5"],
      ["11:00", "PHILADELPHIA", "7"]
    ]
  },
  "config": {
    "rowCount": 3,
    "columnCount": 3
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
