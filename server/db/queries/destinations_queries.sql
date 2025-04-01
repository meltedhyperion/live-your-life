-- name: GetRandomDestinationsForQuestions :many
SELECT id, city, country, clues
FROM destinations
ORDER BY RANDOM()
LIMIT 5;

-- name: GetRandomDestinations :many
SELECT id, city, country
FROM destinations
WHERE id NOT IN (SELECT unnest($1::int[]))
ORDER BY RANDOM()
LIMIT 15;

-- name: GetDestinationByID :one
SELECT city, country, fun_facts, trivia
FROM destinations
WHERE id = $1;


-- -- name: GetRandomDestinationsForSessionQuestions :many
-- SELECT id, city, country, clues
-- FROM destinations
-- ORDER BY RANDOM()
-- LIMIT $1;

-- name: GetRandomDestinationForSessions :many
SELECT id, city, country
FROM destinations
WHERE id NOT IN (SELECT unnest($1::int[]))
ORDER BY RANDOM()
LIMIT $2;

-- name: GetRandomDestinationsForSessionQuestions :many
SELECT id FROM destinations
ORDER BY RANDOM()
LIMIT 10;
