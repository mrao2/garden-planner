-- +goose Up
CREATE TABLE IF NOT EXISTS planters (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  garden_id UUID NOT NULL REFERENCES gardens(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  notes TEXT NOT NULL DEFAULT '',
  type TEXT NOT NULL DEFAULT '',
  planter_type TEXT NOT NULL DEFAULT '',
  shape TEXT NOT NULL DEFAULT 'rectangle',
  sun_level TEXT NOT NULL DEFAULT '',
  watering_type TEXT NOT NULL DEFAULT '',
  width DOUBLE PRECISION NOT NULL DEFAULT 0,
  length DOUBLE PRECISION,
  height DOUBLE PRECISION NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_planters_garden_id ON planters(garden_id);

-- +goose Down
DROP TABLE IF EXISTS planters;
