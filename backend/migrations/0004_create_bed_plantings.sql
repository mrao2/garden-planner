-- +goose Up
CREATE TABLE IF NOT EXISTS bed_plantings (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  bed_id UUID NOT NULL REFERENCES beds(id) ON DELETE CASCADE,
  plant_id UUID NOT NULL REFERENCES plants(id) ON DELETE RESTRICT,
  start_date DATE NOT NULL,
  end_date DATE,
  notes TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_bed_plantings_bed_id ON bed_plantings(bed_id);
CREATE INDEX IF NOT EXISTS idx_bed_plantings_plant_id ON bed_plantings(plant_id);

-- +goose Down
DROP TABLE IF EXISTS bed_plantings;
