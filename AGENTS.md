# Repository Guidelines

## Goal
Build a local-first garden planner:
- `frontend/` SvelteKit static UI.
- `backend/` Go JSON API.
- Local stack via Docker Compose (Caddy + API + Postgres + MinIO).
- MVP focus: beds (SVG rectangles), plants, yearly plan, tasks, journal (text).

## Project Structure & Module Organization
- `backend/` Go API service. Entry point: `backend/api/main.go`. HTTP layer in `backend/internal/httpapi/`; config in `backend/internal/config/`.
- `frontend/` SvelteKit UI (static build). Components should live under `src/lib/components/`.
- `dev/` local tooling config (e.g., `dev/caddyfile`).
- `docker-compose.yml` for the local stack (Caddy + API + Postgres + MinIO when defined).

## Build, Test, and Development Commands
This repo is still being scaffolded, but these are the intended commands:
- `docker compose up --build` runs the full local stack.
- `go run ./backend/api` runs the API directly for quick iteration.
- `go test ./backend/...` runs backend tests (once added).

## Coding Style & Naming Conventions
- Go code must be `gofmt` formatted and follow idiomatic naming (CamelCase for exported, lowerCamelCase for locals).
- Keep backend packages lowercase and focused (e.g., `httpapi`, `config`).
- Frontend should use Tailwind tokens rather than hardcoded hex colors, and avoid heavy UI frameworks unless requested.

## Testing Guidelines
- Prefer Go’s standard `testing` package with `*_test.go` files near the code under test.
- For new endpoints, add minimal handler tests when feasible.
- Keep `/api/health` working at all times.

## Commit & Pull Request Guidelines
- No formal commit style yet; use clear, present-tense messages like `Add router skeleton`.
- PRs should include a short summary, related issue links, and screenshots or logs when behavior is user-visible.

## Configuration & Security Notes
- Never commit secrets. Use environment variables or `.env` files and add them to `.gitignore`.
- Centralize backend configuration in `backend/internal/config/`.
