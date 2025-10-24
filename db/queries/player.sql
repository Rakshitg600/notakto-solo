-- name: GetPlayerById :one
SELECT * FROM Player WHERE uid = $1;

