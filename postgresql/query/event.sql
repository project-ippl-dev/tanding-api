-- name: EventCreate :one
INSERT INTO events(user_id, type, name, description, prize_pool, location, province, city, thumbnail, start_date, end_date, deadline,
                   sport_id, rules, proposal_link, quota, open, remark, status, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, NOW())
RETURNING id;

-- name: EventFetchAll :many
SELECT e.id, e.user_id, u.name as user_name, e.type, e.name, e.description, e.prize_pool, e.location, e.thumbnail, e.start_date, e.end_date,
       e.deadline, e.sport_id, e.remark, e.open, s.name as sport_name, e.rules, e.proposal_link, e.status, e.quota,
       e.province, e.city, u.photo as user_image
FROM events as e
INNER JOIN users as u ON e.user_id = u.id
INNER JOIN sports as s ON e.sport_id = s.id
WHERE e.deleted_at IS NULL
ORDER BY e.created_at DESC
LIMIT $1 OFFSET $2;

-- name: EventCountAll :one
SELECT COUNT(id) FROM events WHERE deleted_at IS NULL LIMIT 1;

-- name: EventUpdate :exec
UPDATE events SET name = $1,
                  type = $2,
                  description = $3,
                  prize_pool = $4,
                  location = $5,
                  province = $6,
                  city = $7,
                  thumbnail = $8,
                  start_date = $9,
                  end_date = $10,
                  deadline = $11,
                  sport_id = $12,
                  rules = $13,
                  proposal_link = $14,
                  quota = $15,
                  open = $16,
                  updated_at = NOW()
WHERE id = $17;

-- name: EventFetchOne :one
SELECT e.id, e.user_id, u.name as user_name, u.photo as user_image, e.type, e.name, e.description, e.prize_pool, e.location, e.thumbnail, e.start_date, e.end_date,
    e.deadline, e.sport_id, s.name as sport_name, e.rules, e.proposal_link, e.status, e.quota, e.remark, e.open,
       e.province, e.city, e.is_generate
FROM events as e
    INNER JOIN users as u ON e.user_id = u.id
    INNER JOIN sports as s ON e.sport_id = s.id
WHERE e.deleted_at IS NULL
AND e.id = $1 LIMIT 1;

-- name: EventCheckOne :one
SELECT id, name, deadline, remark, is_generate FROM events WHERE deleted_at IS NULL AND id = $1;

-- name: EventDelete :exec
UPDATE events SET deleted_at = NOW(),
                  updated_at = NOW()
WHERE id = $1;

-- name: EventFetchByUserID :many
SELECT e.id, e.user_id, u.name as user_name, u.photo as user_image, e.type, e.name, e.description, e.prize_pool, e.location, e.thumbnail, e.start_date, e.end_date,
       e.deadline, e.sport_id, e.remark, e.open, s.name as sport_name, e.rules, e.proposal_link, e.status, e.quota,
       e.province, e.city
FROM events as e
         INNER JOIN event_privileges as ep ON ep.event_id = e.id AND ep.user_id = $1
         INNER JOIN users as u ON e.user_id = u.id
         INNER JOIN sports as s ON e.sport_id = s.id
WHERE e.deleted_at IS NULL
ORDER BY e.created_at DESC
LIMIT $2 OFFSET $3;

-- name: EventCountByUserID :one
SELECT COUNT(id) FROM events WHERE deleted_at IS NULL AND user_id = $1 LIMIT 1;

-- name: EventUpdateStatus :exec
UPDATE events SET status = $1,
                  remark = $2,
                  updated_at = NOW()
WHERE id = $3;

-- name: EventFetchOneByIDAndUserID :one
SELECT id FROM events WHERE id = $1 AND user_id = $2 LIMIT 1;

-- name: EventUpdateRemark :exec
UPDATE events SET remark = $1,
                  updated_at = NOW()
WHERE id = $2;

-- name: EventFetchByRemarkSoon :many
SELECT id, open FROM events WHERE remark = 'soon' AND open <= NOW()  AND status = TRUE;

-- name: EventFetchByRemarkOpen :many
SELECT id, deadline FROM events WHERE remark = 'open' AND deadline <= NOW()  AND status = TRUE;

-- name: EventFetchByRemarkClose :many
SELECT id, start_date FROM events WHERE remark = 'closed' AND is_generate = true AND start_date <= NOW()  AND status = TRUE;

-- name: EventFetchByRemarkOngoing :many
SELECT id, end_date FROM events WHERE remark = 'ongoing' AND end_date < NOW()  AND status = TRUE;

-- name: EventFetchForCart :one
SELECT u.name as event_owner, e.name as event_name, s.name as sport_name, e.thumbnail, e.deadline
FROM events as e
INNER JOIN users as u ON u.id = e.user_id
INNER JOIN sports as s ON s.id = e.sport_id
WHERE e.id = $1;

-- name: EventUpdateIsGenerate :exec
UPDATE events SET is_generate = $1,
                  updated_at = NOW()
WHERE id = $2;


-- name: EventCountByRemark :one
SELECT COUNT(id) FROM events WHERE remark = $1;

-- name: EventCountCanceled :one
SELECT COUNT(id) FROM events WHERE remark = 'rejected';

-- name: EventFetchMostSport :one
SELECT e.sport_id, s.name, COUNT(e.id) as total FROM events as e
INNER JOIN sports as s ON s.id = e.sport_id
GROUP BY e.sport_id, s.name, e.remark
HAVING e.remark = 'done'
ORDER BY total DESC;

-- name: EventFetchMostParticipant :one
SELECT e.id, e.name, COUNT(er.id) as total FROM events as e
INNER JOIN event_registrations as er ON er.event_id = e.id
GROUP BY e.id, e.name, e.remark
HAVING e.remark = 'done'
ORDER BY total DESC;

-- name: EventFetchAllMostSport :many
SELECT s.name, COUNT(e.id) as total FROM events as e
INNER JOIN sports as s ON s.id = e.sport_id
GROUP BY e.sport_id, s.name, e.remark
HAVING e.remark = 'done'
ORDER BY total DESC
LIMIT 10;

-- name: EventFetchAllMostParticipant :many
SELECT e.name, COUNT(er.id) as total FROM events as e
INNER JOIN event_registrations as er ON er.event_id = e.id
GROUP BY e.id, e.name, e.remark
HAVING e.remark = 'done'
ORDER BY total DESC
LIMIT 10;

-- name: EventFetchOneInfiniteByID :one
SELECT e.id, e.user_id, u.name as user_name, e.type, e.name, e.description, e.prize_pool, e.location, e.thumbnail, e.start_date, e.end_date,
       e.deadline, e.sport_id, s.name as sport_name, e.order_number, e.quota, e.remark, e.open, e.province, e.city, u.photo as user_image
FROM events as e
         INNER JOIN users as u ON e.user_id = u.id
         INNER JOIN sports as s ON e.sport_id = s.id
WHERE e.deleted_at IS NULL AND e.remark != 'unconfirmed' AND e.status = TRUE AND e.id = $1;

-- name: EventCountInfinite :one
SELECT COUNT(id) FROM events WHERE remark != 'unconfirmed' AND remark != 'rejected';