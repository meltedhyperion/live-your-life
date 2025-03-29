-- name: GetPlayerById :one
SELECT * FROM players WHERE id = $1;

-- name: UpdatePlayerScore :exec
UPDATE players
SET
    correct_answers = $1,
    total_attempts = $2,
    score = $3,
    updated_at = $4
WHERE id = $5;


-- name: GetLeaderboardDetails :many
SELECT name, avatar, correct_answers, total_attempts, score 
FROM players
WHERE id IN (SELECT unnest($1::uuid[]));

-- name: GetLeaderboardForFriends :many
SELECT p.name, p.avatar, p.correct_answers, p.total_attempts, p.score 
FROM players p
WHERE p.id = $1

UNION

SELECT p.name, p.avatar, p.correct_answers, p.total_attempts, p.score 
FROM players p
JOIN friends f ON f.player2_id = p.id
WHERE f.player1_id = $1;

-- name: CreateNewPlayer :exec
INSERT INTO players (
    id, 
    avatar, 
    name, 
    correct_answers, 
    total_attempts, 
    score, 
    created_at, 
    updated_at
)
VALUES (
    $1, 
    $2, 
    $3, 
    0, 
    0, 
    0.0, 
    CURRENT_TIMESTAMP, 
    CURRENT_TIMESTAMP
);