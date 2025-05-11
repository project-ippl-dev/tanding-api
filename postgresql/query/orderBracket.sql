-- name: OrderBracketCreate :exec
INSERT INTO order_brackets(event_id, class_event_id, club_id, event_registration_id, updated_at)
VALUES ($1, $2, $3, $4, NOW());

-- name: OrderBracketFetchByClassEventID :many
SELECT ob.id, ob.rank, ob.order_by,
       array(SELECT u.name FROM event_participants as ep
           INNER JOIN users as u ON u.id = ep.user_id
           WHERE ep.event_registration_id = ob.event_registration_id
        ) as participants, c.name as club_name
FROM order_brackets as ob
INNER JOIN clubs as c ON c.id = ob.club_id
WHERE ob.class_event_id = $1 ORDER BY ob.order_by;

-- name: OrderBracketUpdateOrderBy :exec
UPDATE order_brackets SET order_by = $1,
                          updated_at = NOW()
WHERE event_registration_id = $2;

-- name: OrderBracketDeleteByClassEventID :exec
DELETE FROM order_brackets WHERE class_event_id = $1;

-- name: OrderBracketCheckOne :one
SELECT ob.id, ob.order_by, ce.score_lock FROM order_brackets as ob
INNER JOIN class_events as ce ON ce.id = ob.class_event_id
WHERE ob.id = $1;

-- name: OrderBracketFetchByClassEventIDAndOrderByScore :many
SELECT ob.id, ob.event_registration_id, ob.class_event_id, ob.club_id, os.total
FROM order_brackets as ob
INNER JOIN order_scores as os ON os.order_bracket_id = ob.id
WHERE ob.class_event_id = $1
ORDER BY os.total DESC;

-- name: OrderBracketUpdateRank :exec
UPDATE order_brackets SET rank = $1,
                          updated_at = NOW()
WHERE id = $2;