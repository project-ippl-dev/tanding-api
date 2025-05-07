-- name: ClubCreate :one
INSERT INTO clubs(name, logo, phone, short_name, user_id, updated_at) VALUES ($1, $2, $3, $4, $5, NOW()) RETURNING id;

-- name: ClubCheckOne :one
SELECT id FROM clubs WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL LIMIT 1;

-- name: ClubCheckOneWithoutUserID :one
SELECT id, name FROM clubs WHERE id = $1 AND deleted_at IS NULL LIMIT 1;

-- name: ClubUpdate :exec
UPDATE clubs SET name = $1,
                 logo = $2,
                 phone = $3,
                 short_name = $4,
                 updated_at = NOW()
WHERE id = $5 AND deleted_at IS NULL;

-- name: ClubDelete :exec
UPDATE clubs SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL;


