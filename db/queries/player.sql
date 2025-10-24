-- name: GetPlayerById :one
SELECT * FROM Player WHERE uid = $1;

-- name: CreatePlayer :exec
INSERT INTO Player (uid, email, name, profile_pic)
VALUES ($1, $2, $3, $4);