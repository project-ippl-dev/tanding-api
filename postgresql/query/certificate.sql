-- name: CertificateCreate :exec
INSERT INTO certificates(id, user_id, event_id, reward_as, class_event_id, updated_at)
VALUES ($1, $2, $3, $4, $5, NOW());

-- name: CertificateFetchOne :one
SELECT c.id, u.name, c.event_id, c.user_id, e.name as event_name, c.reward_as, c.created_at FROM certificates as c
INNER JOIN users as u ON u.id = c.user_id
INNER JOIN events as e ON e.id = c.event_id
WHERE c.id = $1;

-- name: CertificateCountAllByUserID :one
SELECT COUNT(id) FROM certificates WHERE user_id = $1;

-- name: CertificateFetchAllByUserID :many
SELECT c.id, u.name, e.name as event_name, e.thumbnail, c.reward_as, c.created_at FROM certificates as c
INNER JOIN users as u ON u.id = c.user_id
INNER JOIN events as e ON e.id = c.event_id
WHERE c.user_id = $1
ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: CertificateCountAll :one
SELECT COUNT(id) FROM certificates;
