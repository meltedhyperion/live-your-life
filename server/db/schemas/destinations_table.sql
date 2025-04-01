CREATE TABLE destinations (
    id SERIAL PRIMARY KEY,
    city TEXT NOT NULL,
    country TEXT NOT NULL,
    clues JSONB NOT NULL,
    fun_facts JSONB NOT NULL,
    trivia JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT unique_city_country UNIQUE (city, country)
);

ALTER TABLE destinations ENABLE ROW LEVEL SECURITY;


-- Allow public read access to destinations
CREATE POLICY "Allow public read access to destinations"
  ON destinations
  FOR SELECT
  USING (true);

CREATE POLICY "Admin can insert destinations"
  ON destinations
  FOR INSERT
  WITH CHECK (auth.role() = 'admin');

CREATE POLICY "Admin can update destinations"
  ON destinations
  FOR UPDATE
  USING (auth.role() = 'admin');

CREATE POLICY "Admin can delete destinations"
  ON destinations
  FOR DELETE
  USING (auth.role() = 'admin');