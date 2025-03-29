-- name: GetFriendsIdListOfPlayerByID :many
SELECT player2_id FROM friends WHERE player1_id = $1;

-- name: AddFriend :exec
INSERT INTO friends (player1_id, player2_id)
VALUES ($1, $2),
       ($2, $1);