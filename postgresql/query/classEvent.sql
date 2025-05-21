-- name: ClassEventAssign :exec
INSERT INTO class_events (class_id, event_id, price, updated_at) VALUES ($1, $2, $3, NOW());

-- name: ClassEventDetach :exec
DELETE FROM class_events WHERE id = $1 AND event_id = $2;

-- name: ClassEventFetchAll :many
SELECT ce.id, c.name as class_name FROM class_events as ce
    INNER JOIN classes as c ON c.id = ce.class_id
    INNER JOIN events as e ON e.id = ce.event_id
WHERE e.id = $1 AND e.user_id = $2 ORDER BY ce.created_at DESC;

-- name: ClassEventCheckOne :one
SELECT ce.id FROM class_events as ce INNER JOIN events as e ON e.id = ce.event_id
WHERE ce.id = $1 AND ce.event_id = $2 AND e.user_id = $3 LIMIT 1;

-- name: ClassEventUpdate :exec
UPDATE class_events SET price = $1,
                        updated_at = NOW()
WHERE id = $2;

-- name: ClassEventFetchByEventID :many
SELECT ce.id, ce.class_id, c.name as class_name, ce.price, c.match_type, cpr.name as class_rule_name,
       cpr.male as class_rule_male, cpr.female as class_rule_female, cpr.total as class_rule_total
FROM class_events as ce
INNER JOIN classes as c ON c.id = ce.class_id
INNER JOIN class_competition_rules as cpr ON cpr.id = c.class_competition_rule_id
INNER JOIN events as e ON e.id = ce.event_id
WHERE e.id = $1 ORDER BY ce.created_at DESC;

-- name: ClassEventFetchOne :one
SELECT ce.id, c.match_type, ce.bracket_generate, ce.bracket_lock, ce.score_lock, ce.price, ccr.male as rule_male,
       ccr.female as rule_female, ccr.total as rule_total, ce.event_id, ce.match_index
FROM class_events as ce
INNER JOIN classes as c ON c.id = ce.class_id
INNER JOIN class_competition_rules as ccr ON ccr.id = c.class_competition_rule_id
WHERE ce.id = $1;

-- name: ClassEventUpdateBracketGenerate :exec
UPDATE class_events SET bracket_generate = $1,
                        match_index = $2,
                        updated_at = NOW()
WHERE id = $3;

-- name: ClassEventFetchByPaymentID :many
SELECT DISTINCT ce.id, ce.price, c.name FROM class_events as ce
INNER JOIN classes as c ON c.id = ce.class_id
INNER JOIN event_registrations as er ON er.class_event_id = ce.id
INNER JOIN event_payment_receipts as epr ON epr.id = er.event_payment_receipt_id
WHERE epr.id = $1;

-- name: ClassEventWithParticipantsByPaymentID :many
SELECT DISTINCT ce.id, ce.price, c.name, array(SELECT ) FROM class_events as ce
INNER JOIN classes as c ON c.id = ce.class_id
INNER JOIN event_registrations as er ON er.class_event_id = ce.id
INNER JOIN event_payment_receipts as epr ON epr.id = er.event_payment_receipt_id
WHERE epr.id = $1;

-- name: ClassEventUpdateBracketLock :exec
UPDATE class_events SET bracket_lock = $1,
                        updated_at = NOW()
WHERE id = $2;

-- name: ClassEventFetchAndGroupByEventID :many
SELECT ce.id, c.name, ce.match_index, count(eb.id) as total_match,
       count(eb.status) FILTER (WHERE eb.status = 'bye') AS total_bye
FROM class_events as ce
INNER JOIN classes as c ON c.id = ce.class_id
INNER JOIN event_brackets as eb ON eb.class_event_id = ce.id
WHERE ce.event_id = $1
GROUP BY ce.id, c.name ORDER BY ce.match_index DESC, total_match DESC, total_bye;

-- name: ClassEventFetchSingleByEventIDAndLockStatus :many
SELECT ce.id FROM class_events as ce
INNER JOIN classes as c ON c.id = ce.class_id
WHERE c.match_type = 'single' AND ce.bracket_lock = false AND ce.event_id = $1;

-- name: ClassEventFetchLastMatchIndexByEventID :one
SELECT MAX(match_index)::int as last_match_index FROM class_events WHERE event_id = $1;

-- name: ClassEventUpdateScoreLock :exec
UPDATE class_events SET score_lock = $1,
                        updated_at = NOW()
WHERE id = $2;

-- name: ClassEventFetchByEventIDAndScoreLockTrue :many
SELECT ce.id, ce.class_id, c.name as class_name, ce.price, c.match_type, cpr.name as class_rule_name,
       cpr.male as class_rule_male, cpr.female as class_rule_female, cpr.total as class_rule_total, ce.match_index
FROM class_events as ce
         INNER JOIN classes as c ON c.id = ce.class_id
         INNER JOIN class_competition_rules as cpr ON cpr.id = c.class_competition_rule_id
         INNER JOIN events as e ON e.id = ce.event_id
WHERE e.id = $1 AND ce.score_lock = TRUE ORDER BY ce.created_at DESC;