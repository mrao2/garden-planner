-- Demo seed data for local development

-- Gardens
WITH g1 AS (
  INSERT INTO gardens (name, notes)
  VALUES ('Backyard Beds', 'Full sun, raised beds')
  RETURNING id
), g2 AS (
  INSERT INTO gardens (name, notes)
  VALUES ('Kitchen Patio', 'Containers + full sun')
  RETURNING id
)
-- Beds for garden 1 (4 beds)
INSERT INTO beds (garden_id, name, notes, type, shape, sun_level, watering_type, width, length, height)
SELECT g1.id,
       'Bed ' || chr(64 + s) AS name,
       'Full sun' AS notes,
       'raised' AS type,
       'rectangle' AS shape,
       'full-sun' AS sun_level,
       'hand-watered' AS watering_type,
       100 AS width,
       60 AS length,
       18 AS height
FROM g1, generate_series(1, 4) AS s;

-- Beds for garden 2 (12 beds)
WITH g2 AS (
  SELECT id FROM gardens WHERE name = 'Kitchen Patio' ORDER BY created_at DESC LIMIT 1
)
INSERT INTO beds (garden_id, name, notes, type, shape, sun_level, watering_type, width, length, height)
SELECT g2.id,
       'Patio ' || s AS name,
       'Full sun' AS notes,
       'raised' AS type,
       'rectangle' AS shape,
       'full-sun' AS sun_level,
       'hand-watered' AS watering_type,
       90 AS width,
       50 AS length,
       16 AS height
FROM g2, generate_series(1, 12) AS s;

-- Planters
WITH g2 AS (
  SELECT id FROM gardens WHERE name = 'Kitchen Patio' ORDER BY created_at DESC LIMIT 1
)
INSERT INTO planters (garden_id, name, notes, type, planter_type, shape, sun_level, watering_type, width, length, height)
SELECT g2.id,
       'Hanging Basket ' || s AS name,
       'Herb mix' AS notes,
       'basket' AS type,
       'hanging' AS planter_type,
       'circle' AS shape,
       'partial-sun' AS sun_level,
       'hand-watered' AS watering_type,
       18 AS width,
       NULL AS length,
       12 AS height
FROM g2, generate_series(1, 3) AS s;

-- Plants
INSERT INTO plants (name, variety, season, spacing, notes)
VALUES
  ('Tomato', 'San Marzano', 'Warm', '18-24 in', 'Stake early.'),
  ('Tomato', 'Cherokee Purple', 'Warm', '18-24 in', 'Prune lower leaves.'),
  ('Pepper', 'Sweet Bell', 'Warm', '18 in', 'Needs heat.'),
  ('Basil', 'Genovese', 'Warm', '10-12 in', 'Pinch tips.'),
  ('Lettuce', 'Butterhead', 'Cool', '8-10 in', 'Succession sow.'),
  ('Carrot', 'Nantes', 'Cool', '2-3 in', 'Thin twice.'),
  ('Cucumber', 'Marketmore', 'Warm', '12 in', 'Trellis recommended.'),
  ('Zucchini', 'Black Beauty', 'Warm', '24-36 in', 'Give space.' );
