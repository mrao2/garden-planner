-- +goose Up
ALTER TABLE plants
  ADD COLUMN IF NOT EXISTS when_to_sow TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE plants
  DROP COLUMN IF EXISTS when_to_sow;
