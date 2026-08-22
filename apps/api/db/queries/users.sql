-- name: GetUserByID :one
SELECT id, display_name, status
FROM users
WHERE id = $1;
