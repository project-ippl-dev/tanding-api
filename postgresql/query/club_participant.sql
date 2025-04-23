-- name: ClubParticipantCreate :exec
INSERT INTO club_participants(club_id, user_id, club_approval, user_approval, updated_at) VALUES ($1, $2, TRUE, TRUE, NOW());
