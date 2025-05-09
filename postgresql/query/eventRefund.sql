-- name: EventRefundCreate :exec
INSERT INTO event_refunds(event_id, event_payment_receipt_id, note, admin_id, updated_at)
VALUES ($1, $2, $3, $4, NOW());

-- name: EventRefundUpdate :exec
UPDATE event_refunds SET status = $1,
                         admin_id = $2,
                         updated_at = NOW()
WHERE id = $3;

-- name: EventRefundFetchOne :one
SELECT id, event_payment_receipt_id, status FROM event_refunds WHERE id = $1;