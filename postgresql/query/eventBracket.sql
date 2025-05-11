-- name: EventBracketCreate :one
INSERT INTO event_brackets(event_id, class_event_id, event_turn, match_index, match_order, next_match_id, is_active,
                           status, updated_at)
VALUES($1, $2, $3, $4, $5, $6, $7, $8, NOW())
RETURNING id;

-- name: EventBracketUpdateStatus :exec
UPDATE event_brackets SET status = $1,
                          is_active = $2,
                          updated_at = NOW()
WHERE id = $3;

-- name: EventBracketUpdateNextMatch :exec
UPDATE event_brackets SET next_match_id = $1,
                          updated_at = NOW()
WHERE class_event_id = $2 AND match_index = $3 AND match_order = $4;

-- name: EventBracketFetchByMatchIndexAndClassID :many
SELECT id, match_index, match_order FROM event_brackets
WHERE match_index = $1 AND class_event_id = $2 ORDER BY match_index, match_order;

-- name: EventBracketFetchOneByMatchIndexAndClassID :one
SELECT id, match_index, match_order FROM event_brackets
WHERE match_index = $1 AND class_event_id = $2 ORDER BY match_index, match_order;

-- name: EventBracketFetchByClassEventID :many
SELECT eb.id, eb.match_index, eb.match_order, eb.status, eb.event_turn, eb.is_active, COALESCE(ec.home_total,0),
       COALESCE(ec.away_total,0)
FROM event_brackets as eb
LEFT JOIN event_scores as ec ON ec.event_bracket_id = eb.id
WHERE eb.class_event_id = $1 AND eb.match_index = $2 ORDER BY eb.match_index DESC,eb.match_order;

-- name: EventBracketUpdateEventTurnByID :exec
UPDATE event_brackets SET event_turn = $1,
                          updated_at = NOW()
WHERE id = $2;

-- name: EventBracketFetchByClassEventStatusBattleAndMatchIndex :many
SELECT id, match_index, status, event_turn FROM event_brackets
WHERE status = 'battle' AND class_event_id = $1 AND match_index = $2 ORDER BY match_order;

-- name: EventBracketCheckOne :one
SELECT eb.id, eb.next_match_id, eb.match_order, eb.is_active, ce.score_lock FROM event_brackets as eb
INNER JOIN class_events as ce ON ce.id = eb.class_event_id
WHERE eb.id = $1;

-- name: EventBracketFetchByClassEventStatusByeAndMatchIndex :many
SELECT id, match_index, match_order, status, event_turn, next_match_id FROM event_brackets
WHERE status = 'bye' AND class_event_id = $1 AND match_index = $2 ORDER BY match_order;

-- name: EventBracketFetchByClassEventIDAndMatchIndex :many
SELECT eb.id, es.home_total, es.away_total FROM event_brackets AS eb
INNER JOIN event_scores AS es ON es.event_bracket_id = eb.id
WHERE eb.class_event_id = $1 AND eb.match_index = $2;

