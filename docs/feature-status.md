# Feature Status

Status snapshot for the garden-planner MVP. Focus areas: beds (SVG rectangles), plants, yearly plan, tasks, journal.

## Backend (Go API + DB)
Done
- /api/health endpoint. backend/internal/httpapi/router.go
- Gardens CRUD (/api/gardens, /api/gardens/{id}) with validation. backend/internal/httpapi/handlers.go
- Beds CRUD (/api/gardens/{gardenId}/beds, /api/beds/{bedId}) with geometry fields. backend/internal/httpapi/handlers.go
- Plants CRUD (/api/plants, /api/plants/{id}). backend/internal/httpapi/handlers.go
- DB migrations for gardens, beds, plants. backend/migrations/*
- DB connection + migrations on startup. backend/internal/db/*, backend/internal/config/config.go
- Seed data for gardens/beds/plants. dev/seed.sql

Needs work
- Tasks/planner data model + endpoints (not present).
- Journal data model + endpoints (not present).
- Yearly plan model + endpoints (not present).
- Bed ↔ plant assignment / planting plan (not present).
- API for media/attachments (MinIO is configured but unused).

## Frontend (SvelteKit)
Done
- Global layout/navigation + active garden selector. frontend/src/routes/+layout.svelte
- Dashboard with garden CRUD modal, beds snapshot, bed designer. frontend/src/routes/+page.svelte
- Garden store for load/add/edit/delete and selected garden. frontend/src/lib/stores/garden.js
- Beds page with CRUD + geometry edit and BedDesigner commit events. frontend/src/routes/beds/+page.svelte
- Bed designer drag/resize SVG. frontend/src/lib/components/BedDesigner.svelte
- Plants page with list + create + search query. frontend/src/routes/plants/+page.svelte
- Planner page UI with static, hardcoded tasks. frontend/src/routes/planner/+page.svelte
- Journal page UI with local-only entries. frontend/src/routes/journal/+page.svelte
- API client for gardens/beds/plants. frontend/src/lib/api.ts

Needs work
- Planner tasks should come from API (currently hardcoded).
- Journal entries should persist (currently in-memory only).
- Yearly plan UI/data flow is not implemented.
- Tasks CRUD UI and API integration not implemented.
- Bed–plant planning UI and data linkage not implemented.

## Infra / Local Stack
Done
- Docker Compose for Caddy + API + Postgres + MinIO. docker-compose.yml, dev/Caddyfile

Needs work
- MinIO integration in backend (currently unused).
