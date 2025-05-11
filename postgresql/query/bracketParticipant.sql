-- name: BracketParticipantCreate :exec
INSERT INTO bracket_participants(event_registration_id, event_bracket_id, type, is_bye, updated_at)
VALUES ($1, $2, $3, $4, NOW());

-- name: BracketParticipantFetchByEventBracketID :many
SELECT bp.id, er.club_id, c.name as club_name, bp.type, bp.event_registration_id FROM bracket_participants as bp
INNER JOIN event_registrations as er ON er.id = bp.event_registration_id
INNER JOIN clubs as c ON c.id  = er.club_id
WHERE bp.event_bracket_id = $1
ORDER BY CASE bp.type
             WHEN 'home' THEN 1
             WHEN 'away' THEN 2
             END
;

-- name: BracketParticipantCheckOneByEventBracketIDAndType :one
SELECT id FROM bracket_participants WHERE event_bracket_id = $1 AND type = $2;

-- name: BracketParticipantUpdate :exec
UPDATE bracket_participants SET event_registration_id = $1,
                                updated_at = NOW()
WHERE id = $2;

-- name: BracketParticipantUpdateByParticipantType :exec
UPDATE bracket_participants SET event_registration_id = $1,
                                updated_at = NOW()
WHERE event_bracket_id = $2 AND type = $3;

-- name: BracketParticipantFetchByBracketIDAndType :one
SELECT bp.event_registration_id, er.club_id FROM bracket_participants as bp
INNER JOIN event_registrations as er ON er.id = bp.event_registration_id
WHERE bp.event_bracket_id = $1 AND bp.type = $2;