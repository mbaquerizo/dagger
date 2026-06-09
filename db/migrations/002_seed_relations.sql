INSERT INTO relations (name)
VALUES
  ('blocks'),
  ('blocked_by'),
  ('duplicates'),
  ('duplicated_from'),
  ('relates_to'),
  ('causes'),
  ('caused_by')
ON CONFLICT (name) DO NOTHING;