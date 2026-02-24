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

## Phase 2.5: Display Management

- [ ] Display list view (fetch from `GET /api/v1/displays`)
- [ ] Display selector/switcher in UI
- [ ] Create display form (ID, rows, columns)
- [ ] Edit display content (simple grid editor)
- [ ] Delete display action
- [ ] Make `DisplayPreview` accept dynamic `displayId` prop

> **Architecture Decision**: Phase 2.5 bridges the gap between Phase 2's backend CRUD capabilities and Phase 3's advanced builder features. Currently, the frontend is hardcoded to fetch `/api/v1/displays/demo` with no way to create, switch between, or manage multiple displays. This phase provides basic display management UI, separating concerns: Phase 2.5 handles *which* display and fundamental CRUD operations, while Phase 3 focuses on *how* to configure displays with advanced options (character sets, timing, etc.). This prevents Phase 3 from becoming too large and validates that Phase 2's API contracts work end-to-end.

## Phase 3: Builder

- [ ] Character set selection UI
- [ ] Flip speed/timing controls
- [ ] Advanced configuration options
- [ ] Enhanced preview with real-time updates
- [ ] Configuration presets/templates

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

**Current Phase:** 2.5 (Display Management)

**Out of scope (for now):** Authentication, themes, mobile app, public API
