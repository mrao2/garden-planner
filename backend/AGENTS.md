# AGENTS.md — Backend (Go API)

## Purpose
JSON API for garden planner. Entry point: `backend/api/main.go`.
All routes under `/api/...`.

## Architecture
- Keep code in `backend/internal/...` with clear packages:
  - `internal/httpapi` (router, middleware, handlers)
  - `internal/db` (queries/repositories)
  - `internal/domain` (types/models)
  - `internal/storage` (S3/MinIO client; later)
- Avoid giant files; split by feature area.

## API conventions
- JSON response envelope:
  - success: `{ "data": ... }`
  - error: `{ "error": { "code": "...", "message": "..." } }`
- Use proper HTTP status codes.
- Validate inputs; return consistent error codes.

## Data + migrations
- Use Postgres for structured data.
- Prefer a simple migration tool (goose) when we get there.
- No blobs in Postgres.

## Object storage
- MinIO/S3 via environment vars:
  - `S3_ENDPOINT`, `S3_BUCKET`, `S3_ACCESS_KEY`, `S3_SECRET_KEY`, `S3_USE_PATH_STYLE`
- For uploads, use presigned URLs (client uploads directly).
- MVP does not require media features unless explicitly requested.

## Testing & reliability
- Add at least minimal handler tests for new endpoints when feasible.
- Keep `/api/health` always working.

