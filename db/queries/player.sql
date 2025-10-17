-- name: getPlayerById :one
SELECT * FROM Player WHERE uid = $1;

