-- +goose Up
ALTER TABLE plants
  DROP COLUMN IF EXISTS germination_time;

-- +goose Down
ALTER TABLE plants
  ADD COLUMN IF NOT EXISTS germination_time TEXT NOT NULL DEFAULT '';
