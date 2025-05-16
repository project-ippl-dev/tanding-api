-- name: ClubCertificateCreate :exec
INSERT INTO club_certificates(id, club_id, event_id, reward_as, updated_at)
VALUES($1, $2, $3, $4, NOW());

-- name: ClubCertificateCountAllByClubID :one
SELECT COUNT(id) FROM club_certificates WHERE club_id = $1;

-- name: ClubCertificateFetchAllByUserID :many
SELECT c.id, cl.name, e.name as event_name, e.thumbnail, c.reward_as, c.created_at FROM club_certificates as c
INNER JOIN clubs as cl ON cl.id = c.club_id
INNER JOIN events as e ON e.id = c.event_id
WHERE c.club_id = $1
ORDER BY created_at DESC LIMIT $2 OFFSET $3;