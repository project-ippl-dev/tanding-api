-- name: EventRegistrationCreate :one
INSERT INTO event_registrations(event_id, class_event_id, club_id, event_payment_receipt_id, status, updated_at)
VALUES ($1, $2, $3, $4, $5, NOW())
RETURNING id;

-- name: EventRegistrationUpdate :exec
UPDATE event_registrations SET class_event_id = $1,
                               updated_at = NOW()
WHERE id = $2;

-- name: EventRegistrationFetchOne :one
SELECT id, event_id, class_event_id, club_id, event_payment_receipt_id, status FROM event_registrations
WHERE id = $1;

-- name: EventRegistrationSetReject :exec
UPDATE event_registrations SET status = 'rejected',
                               updated_at = NOW()
WHERE id = $1;

-- name: EventRegistrationUpdatePaymentID :exec
UPDATE event_registrations SET event_payment_receipt_id = $1,
                               status = $2,
                               updated_at = NOW()
WHERE event_payment_receipt_id = $1;

-- name: EventRegistrationUpdateOnePaymentID :exec
UPDATE event_registrations SET event_payment_receipt_id = $1,
                               status = $2,
                               updated_at = NOW()
WHERE id = $3;

-- name: EventRegistrationCountAllByStatusApproved :one
SELECT COUNT(id) FROM event_registrations WHERE event_id = $1 AND status = 'approved';

-- name: EventRegistrationFetchByPaymentID :many
SELECT id, event_id, class_event_id, club_id, status FROM event_registrations
WHERE event_payment_receipt_id = $1;

-- name: EventRegistrationFetchByClassEventID :many
SELECT er.id, er.event_id, er.class_event_id, er.club_id, er.status, c.name as club_name FROM event_registrations as er
INNER JOIN clubs as c ON c.id = er.club_id
WHERE status = 'approved' AND class_event_id = $1;

-- name: EventRegistrationFetchCart :many
SELECT e.id, e.name, e.thumbnail, SUM(ce.price) as total FROM event_registrations as er
INNER JOIN events as e ON e.id = er.event_id
INNER JOIN class_events as ce ON ce.id = er.class_event_id
INNER JOIN clubs as c ON c.id = er.club_id
WHERE er.status = 'pending' AND c.user_id = $1
GROUP BY e.id ORDER BY e.created_at DESC LIMIT $2 OFFSET $3;

-- name: EventRegistrationCountCart :one
SELECT COUNT(*) FROM (SELECT e.id, e.name, e.thumbnail, SUM(ce.price) as total FROM event_registrations as er
    INNER JOIN events as e ON e.id = er.event_id
    INNER JOIN class_events as ce ON ce.id = er.class_event_id
    INNER JOIN clubs as c ON c.id = er.club_id
    WHERE er.status = 'pending' AND c.user_id = $1
    GROUP BY e.id) as carts;

-- name: EventRegistrationFetchPendingCron :many
SELECT er.id FROM event_registrations as er
                      INNER JOIN events as e ON e.id = er.event_id
WHERE e.deadline <= $1 AND e.remark != ANY('{closed,ongoing,done}') AND er.status = 'pending';

-- name: EventRegistrationUpdateStatus :exec
UPDATE event_registrations SET status = $1,
                               updated_at = NOW()
WHERE id = $2;

-- name: EventRegistrationFetchByClassEventIDAndStatus :many
SELECT er.id, c.name as class_name,cl.name as club_name, er.status, er.club_id, er.class_event_id, ce.price FROM event_registrations as er
INNER JOIN class_events as ce ON ce.id = er.class_event_id
INNER JOIN classes as c ON c.id = ce.class_id
INNER JOIN clubs as cl ON cl.id = er.club_id
WHERE er.class_event_id = $1 AND er.status = 'approved';

-- name: EventRegistrationFetchClubByEventID :many
SELECT DISTINCT cl.id, cl.name FROM event_registrations as er
LEFT JOIN ranks as r ON r.event_registration_id = er.id
INNER JOIN clubs as cl ON cl.id = er.club_id
WHERE er.event_id = $1 and er.status = 'approved'
ORDER BY cl.name;

-- name: EventRegistrationCountByEventIDAndClubID :one
SELECT COUNT(id) FROM event_registrations
WHERE event_id = $1 AND club_id = $2 AND status = 'approved';
