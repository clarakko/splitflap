-- SplitFlap Display Database Schema
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
