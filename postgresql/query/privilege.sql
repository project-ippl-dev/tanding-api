-- name: PrivilegeFetchOneByUserID :one
SELECT privileges.id, privileges.name FROM privileges
INNER JOIN privilege_user ON privileges.id = privilege_user.privilege_id
INNER JOIN users ON privilege_user.user_id = users.id
WHERE users.id = $1 AND privileges.name = $2 AND privileges.type = $3 LIMIT 1;

-- name: PrivilegeFetchByUserID :many
SELECT privileges.name FROM privileges
                                               INNER JOIN privilege_user ON privileges.id = privilege_user.privilege_id
                                               INNER JOIN users ON privilege_user.user_id = users.id
WHERE users.id = $1;