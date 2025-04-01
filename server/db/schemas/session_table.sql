CREATE TABLE sessions (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES players(id) NOT NULL,
    destinations INT[] NOT NULL DEFAULT '{}',
    score FLOAT DEFAULT 0 NOT NULL,
    total_attempted INT DEFAULT 0,
    correct INT DEFAULT 0
);