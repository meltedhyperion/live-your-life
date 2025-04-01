-- name: GetUserSessionByID :one
SELECT * FROM sessions WHERE id = $1 AND user_id = $2;


-- name: GetAllUserSessionByID :one
SELECT id FROM sessions WHERE user_id = $1 
ORDER BY id ASC 
LIMIT 1;
-- name: CreateUserSession :exec

INSERT INTO sessions (
    user_id,
    destinations
) VALUES ($1, $2);


-- name: UpdateUserSession :exec
UPDATE sessions
SET
    score = $1,
    total_attempted = $2,
    correct = $3
WHERE id = $4;


