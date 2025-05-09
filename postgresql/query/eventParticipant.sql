-- name: EventParticipantCreate :exec
INSERT INTO event_participants(event_registration_id, user_id, updated_at) VALUES ($1, $2, NOW());

-- name: EventParticipantFetchByRegistrationID :many
SELECT ep.id, ep.event_registration_id, ep.user_id, u.name, c.name as class_name FROM event_participants as ep
INNER JOIN event_registrations as er ON er.id = ep.event_registration_id
INNER JOIN class_events as ce ON ce.id = er.class_event_id
INNER JOIN classes as c ON c.id = ce.class_id
INNER JOIN users as u ON u.id = ep.user_id
WHERE ep.event_registration_id = $1;

-- name: EventParticipantFetchNameByRegistrationID :many
SELECT u.name FROM event_participants as ep
INNER JOIN users as u ON u.id = ep.user_id
WHERE ep.event_registration_id = $1;

-- name: EventParticipantFetchByEventIDAndClubID :many
SELECT ep.id, ep.event_registration_id, ep.user_id, u.name, c.name as class_name FROM event_participants as ep
INNER JOIN event_registrations as er ON er.id = ep.event_registration_id
INNER JOIN class_events as ce ON ce.id = er.class_event_id
INNER JOIN classes as c ON c.id = ce.class_id
INNER JOIN users as u ON u.id = ep.user_id
WHERE er.event_id = $1 AND er.club_id = $2 AND er.status = 'approved'
ORDER BY u.name;