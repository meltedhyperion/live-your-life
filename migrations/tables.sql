CREATE TABLE players (
    id UUID PRIMARY KEY,              -- Supabase auth user ID
    avatar TEXT,
    name TEXT,
    correct_answers INT DEFAULT 0,
    total_attempts INT DEFAULT 0,
    score FLOAT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE destinations (
    id SERIAL PRIMARY KEY,
    city TEXT NOT NULL,
    country TEXT NOT NULL,
    clues JSONB NOT NULL,
    fun_facts JSONB NOT NULL,
    trivia JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_city_country UNIQUE (city, country)
);

CREATE TABLE friends (
    id SERIAL PRIMARY KEY,
    player1_id UUID NOT NULL,
    player2_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_friendship UNIQUE (player1_id, player2_id),
    CHECK (player1_id <> player2_id),
    FOREIGN KEY (player1_id) REFERENCES players(id) ON DELETE CASCADE,
    FOREIGN KEY (player2_id) REFERENCES players(id) ON DELETE CASCADE
);

ALTER TABLE players ENABLE ROW LEVEL SECURITY;
ALTER TABLE friends ENABLE ROW LEVEL SECURITY;
ALTER TABLE destinations ENABLE ROW LEVEL SECURITY;


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

-- Allow a user to view friendship records if they are either player1 or player2
CREATE POLICY "Players can view their friendships"
  ON friends
  FOR SELECT
  USING (player1_id = auth.uid() OR player2_id = auth.uid());

-- Allow insertion only if the authenticated user is involved in the friendship
CREATE POLICY "Players can insert friendship records"
  ON friends
  FOR INSERT
  WITH CHECK (player1_id = auth.uid() OR player2_id = auth.uid());

CREATE POLICY "Players can update their friendships"
  ON friends
  FOR UPDATE
  USING (player1_id = auth.uid() OR player2_id = auth.uid());

CREATE POLICY "Players can delete their friendships"
  ON friends
  FOR DELETE
  USING (player1_id = auth.uid() OR player2_id = auth.uid());

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
