CREATE TABLE friends (
    id SERIAL PRIMARY KEY,
    player1_id UUID NOT NULL,
    player2_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT unique_friendship UNIQUE (player1_id, player2_id),
    CHECK (player1_id <> player2_id),
    FOREIGN KEY (player1_id) REFERENCES players(id) ON DELETE CASCADE,
    FOREIGN KEY (player2_id) REFERENCES players(id) ON DELETE CASCADE
);

ALTER TABLE friends ENABLE ROW LEVEL SECURITY;

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