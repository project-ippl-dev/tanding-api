-- name: LoginDetailCreate :exec
INSERT INTO login_details(user_id, updated_at) VALUES ($1, NOW());

-- name: LoginDetailCount :one
SELECT COUNT(id) FROM login_details;

-- name: LoginDetailLastRecord :one
SELECT u.name, ld.created_at as last_login FROM login_details as ld
INNER JOIN users as u ON u.id = ld.user_id
ORDER BY ld.created_at DESC;

-- name: LoginDetailFetchAll :many
SELECT u.name, ld.created_at as last_login FROM login_details as ld
INNER JOIN users as u ON u.id = ld.user_id
ORDER BY last_login DESC
LIMIT $1 OFFSET $2;