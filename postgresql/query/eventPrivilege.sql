-- name: EventPrivilegeCreate :exec
INSERT INTO event_privileges (event_id, user_id, role, updated_at) VALUES ($1, $2, $3, NOW());

-- name: EventPrivilegeFetchOne :one
SELECT id, role FROM event_privileges WHERE event_id = $1 AND user_id = $2 LIMIT 1;

-- name: EventPrivilegeFetchAll :many
SELECT ep.id, ep.role, ep.user_id, u.name FROM event_privileges as ep
INNER JOIN users as u ON u.id = ep.user_id
WHERE ep.event_id = $1;

-- name: EventPrivilegeUpdate :exec
UPDATE event_privileges SET role = $1,
                            updated_at = NOW()
WHERE user_id = $2 AND event_id = $3;

-- name: EventPrivilegeDelete :exec
DELETE FROM event_privileges WHERE event_id = $1 AND user_id = $2;

-- name: EventPrivilegeFetchByUserID :many
SELECT ep.id FROM event_privileges as ep
INNER JOIN events as e ON ep.event_id = e.id
WHERE e.deleted_at IS NULL AND ep.user_id = $1;

-- name: EventPrivilegeFetchByEventID :many
SELECT ep.id, ep.user_id, ep.role FROM event_privileges as ep
WHERE ep.event_id = $1;