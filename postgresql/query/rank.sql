-- name: RankCreate :exec
INSERT INTO ranks(club_id, event_id, class_event_id, event_registration_id, sport_id, rank, point, updated_at)
VALUES($1, $2, $3, $4, $5, $6, $7, NOW());

-- name: RankClubFetchByEventID :many
SELECT r.club_id, c.name AS club_name,
       COUNT(*) FILTER (WHERE r.rank = 1) AS rank1,
       COUNT(*) FILTER (WHERE r.rank = 2) AS rank2,
       COUNT(*) FILTER (WHERE r.rank = 3 OR r.rank = 4) AS rank3,
       SUM(r.point) as total_point
FROM ranks as r
         INNER JOIN clubs as c ON c.id = r.club_id
WHERE r.event_id = $1
GROUP BY r.club_id, c.name
ORDER BY rank1 DESC, rank2 DESC, rank3 desc, total_point DESC;

-- name: RankFetchPointByUserID :one
SELECT COALESCE(SUM(r.point),0)::int as total FROM ranks as r
INNER JOIN event_registrations as er ON er.id = r.event_registration_id
INNER JOIN event_participants as ep ON ep.event_registration_id = er.id
WHERE ep.user_id = $1;

-- name: RankFetchPointByClubID :one
SELECT SUM(point) as total FROM ranks WHERE club_id = $1;

-- name: RankFetchAllPointByClubID :many
SELECT u.id, u.name, COALESCE(SUM(r.point),0) as point FROM ranks as r
INNER JOIN event_registrations as er ON er.id = r.event_registration_id
INNER JOIN event_participants as ep ON ep.event_registration_id = er.id
INNER JOIN users as u ON u.id = ep.user_id
GROUP BY u.id, r.club_id
HAVING r.club_id = $1;

-- name: RankFetchPointByClubIDAndUserID :one
SELECT SUM(r.point) as total FROM ranks as r
INNER JOIN event_registrations as er ON er.id = r.event_registration_id
INNER JOIN event_participants as ep ON ep.event_registration_id = er.id
WHERE er.club_id = $1 AND ep.user_id = $2;

-- name: RankCountAllPointUser :one
SELECT COUNT(*) FROM (SELECT u.id, u.name, COALESCE(SUM(r.point),0) as point FROM ranks as r
INNER JOIN event_registrations as er ON er.id = r.event_registration_id
INNER JOIN event_participants as ep ON ep.event_registration_id = er.id
INNER JOIN users as u ON u.id = ep.user_id
GROUP BY u.id, r.club_id
ORDER BY point DESC) as total;