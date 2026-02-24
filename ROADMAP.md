# Roadmap

## Phase 1: MVP

- [X] Backend serves display content via REST endpoint
- [X] Frontend fetches and displays data from API
- [X] Basic flip animation in SolidJS
- [X] Single display with hardcoded rows
- [X] Frontend renders split-flap animation

## Phase 2: Persistence

- [X] Database setup (SQLite)
- [X] Display entity and repository
- [X] CRUD endpoints for displays

> **Architecture Decision**: Using SQLite for both development and production. Provides zero-config setup, sufficient scale for split-flap displays (read-heavy workload), and keeps migration path to Postgres open via repository pattern if needed in Phase 8.

## Phase 3: Builder

- [ ] Configurable rows/columns in UI
- [ ] Character set options
- [ ] Flip speed/timing settings
- [ ] Save display configuration
- [ ] Preview animation

## Phase 4: Embed Component

- [ ] Web component (vanilla JS/TS)
- [ ] Fetches config from API by display ID
- [ ] Lightweight, dependency-free
- [ ] Embed code generator in builder

## Phase 5: Real-time

- [ ] WebSocket support
- [ ] Push updates to connected clients
- [ ] Multiple clients stay in sync

## Phase 6: Data Sources

- [ ] External API integration (weather, transit, etc.)
- [ ] Data source abstraction
- [ ] Caching layer

## Phase 7: Scheduling

- [ ] Scheduled content changes
- [ ] Time-based display rules

## Phase 8: Multi-user

- [ ] User accounts
- [ ] Multiple displays per user
- [ ] Dashboard view

---

**Current Phase:** 1 (MVP)

**Out of scope (for now):** Authentication, themes, mobile app, public API
