-- name: UserCreate :one
INSERT INTO users(name, phone, photo, born_at, born_on, identity_number, gender, about, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
RETURNING id;

-- name: UserFetchOne :one
SELECT users.id, accounts.id as account_id, users.name, users.photo, accounts.username, accounts.type, users.role
FROM users
INNER JOIN accounts ON accounts.user_id = users.id
WHERE accounts.username = $1
ORDER BY array_position(array['manual', 'google', 'facebook']::account_type[] , accounts.type)
LIMIT 1;

-- name: UserFetchOneTypeManual :one
SELECT users.id, accounts.id as account_id, users.name, users.photo, accounts.username, accounts.type, users.role
FROM users
         INNER JOIN accounts ON accounts.user_id = users.id
WHERE accounts.username = $1 AND accounts.type = 'manual'
LIMIT 1;

-- name: UserFetchNameByID :one
SELECT name FROM users WHERE id = $1 LIMIT 1;

-- name: UserFetchPhotoByID :one
SELECT photo FROM users WHERE id = $1 LIMIT 1;

-- name: UserFetchByKeyword :many
SELECT DISTINCT users.id, users.name, accounts.username FROM users
INNER JOIN accounts ON accounts.user_id = users.id
WHERE users.role = 'user' AND users.id != $1 AND users.name LIKE $2 OR accounts.username LIKE $3 LIMIT $4;

-- name: UserFetchBasicInformation :one
SELECT id, name, born_at, born_on, identity_number, phone, photo, gender, about, status, can_participate FROM users
WHERE id = $1;

-- name: UserUpdateBasicInformation :exec
UPDATE users SET name = $1,
                 born_at = $2,
                 born_on = $3,
                 identity_number = $4,
                 phone = $5,
                 photo = $6,
                 gender = $7,
                 about = $8,
                 can_participate = TRUE,
                 updated_at = NOW()
WHERE id = $9;

-- name: UserFetchInID :many
SELECT id, name, gender, can_participate FROM users WHERE id = ANY($1);

-- name: UserCountAll :one
SELECT COUNT(id) FROM users where status = TRUE;

-- name: UserFetchAll :many
SELECT id, name, role, created_at FROM users
ORDER BY created_at DESC LIMIT $1 OFFSET $2;
