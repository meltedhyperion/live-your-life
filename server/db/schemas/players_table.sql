CREATE TABLE players (
    id UUID PRIMARY KEY,              -- Supabase auth user ID
    avatar TEXT NOT NULL,
    name TEXT NOT NULL,
    correct_answers INT DEFAULT 0 NOT NULL,
    total_attempts INT DEFAULT 0 NOT NULL,
    score FLOAT DEFAULT 0 NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);
ALTER TABLE players ENABLE ROW LEVEL SECURITY;

-- Allow a user to insert their own player record (id must equal auth.uid())
CREATE POLICY "Players can insert their own data"
  ON players
  FOR INSERT
  WITH CHECK (id = auth.uid());

-- Allow a user to read their own player record only
CREATE POLICY "Players can read their own data"
  ON players
  FOR SELECT
  USING (id = auth.uid());

-- Allow a user to update their own record only
CREATE POLICY "Players can update their own data"
  ON players
  FOR UPDATE
  USING (id = auth.uid());
