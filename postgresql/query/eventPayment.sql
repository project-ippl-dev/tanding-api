-- name: EventPaymentCreate :one
INSERT INTO event_payment_receipts(unique_number, event_id, user_id ,payment_link, admin_id, total, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, NOW())
RETURNING id;

-- name: EventPaymentCheckOne :one
SELECT ep.id FROM event_payment_receipts as ep
INNER JOIN event_registrations as er ON er.event_payment_receipt_id = ep.id
INNER JOIN clubs as c ON c.id = er.club_id
WHERE er.event_id = $1 AND ep.id = $2;

-- name: EventPaymentUpdate :exec
UPDATE event_payment_receipts SET status = $1,
                                  admin_id = $2,
                                  updated_at = NOW()
WHERE id = $3;

-- name: EventPaymentFetchOne :one
SELECT id, status, total, unique_number, created_at, event_id FROM event_payment_receipts WHERE id = $1 AND user_id = $2;

-- name: EventPaymentFetchOneForAdmin :one
SELECT DISTINCT ep.id, ep.status, ep.total, ep.unique_number, ep.created_at, ep.payment_link,ep.event_id, c.name as club_name,
       u.name as owner, e.name as event_name
FROM event_payment_receipts as ep
INNER JOIN events as e ON e.id = ep.event_id
INNER JOIN event_registrations as er ON er.event_payment_receipt_id = ep.id
INNER JOIN clubs as c ON c.id = er.club_id
INNER JOIN users as u ON u.id = c.user_id
WHERE ep.id = $1;

-- name: EventPaymentSumAll :one
SELECT SUM(c.price) FROM event_registrations AS er
INNER JOIN class_events AS ce ON ce.id = er.class_event_id
INNER JOIN classes as c ON c.id = ce.class_id
WHERE er.status = $1;