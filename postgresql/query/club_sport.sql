-- name: ClubAttachSport :exec
INSERT INTO club_sport(club_id, sport_id, updated_at) VALUES ($1, $2, NOW());
