-- +goose Up
ALTER TABLE plants
  ADD COLUMN IF NOT EXISTS days_to_maturity INTEGER,
  ADD COLUMN IF NOT EXISTS sow_depth TEXT NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS germination_time TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE plants
  DROP COLUMN IF EXISTS days_to_maturity,
  DROP COLUMN IF EXISTS sow_depth,
  DROP COLUMN IF EXISTS germination_time;
