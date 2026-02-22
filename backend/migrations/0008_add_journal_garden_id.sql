-- +goose Up
ALTER TABLE journal_entries
  ADD COLUMN IF NOT EXISTS garden_id UUID;

UPDATE journal_entries
  SET garden_id = (SELECT id FROM gardens ORDER BY created_at ASC LIMIT 1)
  WHERE garden_id IS NULL;

ALTER TABLE journal_entries
  ALTER COLUMN garden_id SET NOT NULL,
  ADD CONSTRAINT journal_entries_garden_id_fkey FOREIGN KEY (garden_id) REFERENCES gardens(id) ON DELETE CASCADE;

-- +goose Down
ALTER TABLE journal_entries
  DROP CONSTRAINT IF EXISTS journal_entries_garden_id_fkey,
  DROP COLUMN IF EXISTS garden_id;
