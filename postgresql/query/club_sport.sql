-- name: ClubAttachSport :exec
INSERT INTO club_sport(club_id, sport_id, updated_at) VALUES ($1, $2, NOW());

-- name: ClubSportFetchByClubID :many
SELECT cs.id, s.id as sport_id ,s.name as sport_name FROM club_sport as cs
INNER JOIN sports as s ON s.id = cs.sport_id
WHERE cs.club_id = $1;