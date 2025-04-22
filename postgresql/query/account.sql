-- name: AccountFetchUserIDByUsername :one
SELECT accounts.user_id FROM accounts INNER JOIN users ON users.id = accounts.user_id
WHERE accounts.username = $1 AND accounts.type = $2 AND users.status = '1' LIMIT 1;

-- name: AccountCreate :one
INSERT INTO accounts(username, password, type, user_id, status, updated_at) VALUES ($1, $2, $3, $4, $5, NOW())
RETURNING id;

-- name: AccountFetchEmailByID :one
SELECT username FROM accounts INNER JOIN users ON users.id = accounts.user_id
WHERE accounts.id = $1 AND accounts.type = $2 AND users.status = '1' LIMIT 1;

-- name: AccountUpdateStatusByID :exec
UPDATE accounts SET status = true,
                    updated_at = NOW()
WHERE id = $1;

-- name: AccountFetchOneByEmail :one
SELECT accounts.id, accounts.username, accounts.password, accounts.user_id,
       users.role, users.photo, users.name, users.can_participate
FROM accounts INNER JOIN users ON users.id = accounts.user_id
WHERE username = $1 AND type = $2 AND users.status = '1' LIMIT 1;

-- name: AccountUpdatePassword :exec
UPDATE accounts SET password = $1,
                    updated_at = NOW()
WHERE id = $2;