# AGENTS.md — Frontend (SvelteKit)

## Purpose
User interface for the garden planner. Served as a static build.

## Design system
### Theme: earthy / garden-inspired
Use a calm, natural palette:
- Primary: deep forest green
- Secondary: clay/terracotta
- Accent: muted mustard
- Neutrals: warm sand, off-white, bark brown
- Avoid neon, pure black (#000), or high-saturation colors.

If you need concrete tokens, use these CSS variables:
- `--c-bg: #F6F1E7` (warm off-white)
- `--c-surface: #EFE6D6` (sand)
- `--c-text: #2B2A27` (bark)
- `--c-primary: #2F5D50` (forest)
- `--c-secondary: #B45A3C` (terracotta)
- `--c-accent: #C2A24C` (mustard)
- `--c-border: #D8CBB5` (warm border)

Keep contrast readable.

## UI conventions
- Use simple, minimal layout; avoid clutter.
- Prefer components in `src/lib/components/`.
- Keep the SVG bed designer isolated (e.g., `BedDesigner.svelte`).
- Use fetch to `/api/*` on the same origin; do not hardcode host/ports.
- Use a small `api.ts` client wrapper in `src/lib/api/`.

## SVG bed editor v1
- Rectangle beds only.
- Interactions: click to select, drag to move, resize via corner handles.
- Grid snapping optional, but keep math simple.
- Bed properties panel: name, sun exposure, slope, notes.

## Build constraints
- Must build as static output (adapter-static) with SPA fallback.
- Do not introduce heavy UI frameworks unless requested.

## Styling
- Use Tailwind CSS for all styling.
- Prefer theme tokens from `tailwind.config.cjs` under `colors.earth.*`.
- Default page background: `bg-earth-bg`, text: `text-earth-text`.
- Use `earth-forest` for primary actions, `earth-terracotta` for highlights, `earth-mustard` sparingly for accents.
- Avoid hardcoded hex colors in components unless adding a new token.
- Keep UI clean and calm; no neon/high-saturation colors.
