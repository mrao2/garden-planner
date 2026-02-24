-- +goose Up
ALTER TABLE plants ADD COLUMN IF NOT EXISTS is_starred BOOLEAN NOT NULL DEFAULT false;

-- +goose Down
ALTER TABLE plants DROP COLUMN IF EXISTS is_starred;
