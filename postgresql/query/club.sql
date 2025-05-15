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

-- name: ClubHandover :exec
UPDATE clubs SET user_id = $1,
                 updated_at = NOW()
WHERE id = $2;

-- name: ClubFetchOne :one
SELECT clubs.id, clubs.name, clubs.logo, clubs.short_name, users.name as owner FROM clubs
INNER JOIN users ON users.id = clubs.user_id
WHERE clubs.id = $1
ORDER BY clubs.created_at DESC
LIMIT 1;

-- name: ClubFetchOwner :one
SELECT id FROM clubs WHERE user_id = $1 LIMIT 1;

-- name: ClubFetchByUserID :many
SELECT c.id, c.name FROM clubs as c
INNER JOIN club_participants as cp ON c.id = cp.club_id AND cp.club_approval = TRUE AND cp.user_approval = TRUE
WHERE cp.user_id = $1;

-- name: ClubCheckByIDAndUserID :one
SELECT id FROM clubs WHERE id = $1 AND user_id = $2;

-- name: ClubFetchOwnerByEventID :many
SELECT c.id FROM clubs as c
INNER JOIN event_registrations as er ON er.club_Id = c.id AND er.event_id = $1
WHERE user_id = $2;

-- name: ClubFetchAllOwner :many
SELECT id, name, logo FROM clubs WHERE user_id = $1;

-- name: ClubFetchOneOwner :one
SELECT id FROM clubs WHERE user_id = $1 AND id = $2 LIMIT 1;