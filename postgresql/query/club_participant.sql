-- name: ClubParticipantCreate :exec
INSERT INTO club_participants(club_id, user_id, club_approval, user_approval, updated_at) VALUES ($1, $2, TRUE, TRUE, NOW());

-- name: ClubParticipantInvite :exec
INSERT INTO club_participants(club_id, user_id, sport_id, club_approval, updated_at) VALUES ($1, $2, $3, TRUE, NOW());

-- name: ClubParticipantJoin :exec
INSERT INTO club_participants(club_id, user_id, sport_id, user_approval, updated_at) VALUES($1, $2, $3, TRUE, NOW());

-- name: ClubParticipantCheckParticipant :one
SELECT id FROM club_participants WHERE club_id = $1 AND user_id = $2 AND club_approval = true AND user_approval = true LIMIT 1;

-- name: ClubParticipantFetchLatestIDJoinApproval :one
SELECT id FROM club_participants WHERE club_id = $1 AND user_approval = TRUE AND club_approval IS NULL ORDER BY id DESC;

-- name: ClubParticipantFetchJoinApproval :many
SELECT cp.id, s.id as sport_id, s.name as sport_name, u.name FROM club_participants as cp
    INNER JOIN users as u ON cp.user_id = u.id
    INNER JOIN sports as s ON cp.sport_id = s.id
WHERE club_id = $1  AND cp.id < $2 AND cp.user_approval = TRUE AND cp.club_approval IS NULL ORDER BY cp.id DESC LIMIT $3;

-- name: ClubParticipantFetchLatestIDInviteApproval :one
SELECT id FROM club_participants WHERE user_id = $1 AND user_approval IS NULL AND club_approval = TRUE ORDER BY id DESC;

-- name: ClubParticipantFetchInviteApproval :many
SELECT cp.id, cp.club_id, s.id as sport_id, s.name as sport_name, c.name FROM club_participants as cp
    INNER JOIN clubs as c ON cp.club_id = c.id
    INNER JOIN sports as s ON cp.sport_id = s.id
WHERE cp.user_id = $1  AND cp.id < $2 AND cp.user_approval IS NULL AND cp.club_approval = TRUE ORDER BY cp.id DESC LIMIT $3;

-- name: ClubParticipantUpdateJoinApproval :exec
UPDATE club_participants SET club_approval = $1,
                             updated_at = NOW()
WHERE id = $2 AND club_id = $3;

-- name: ClubParticipantUpdateInviteApproval :exec
UPDATE club_participants SET user_approval = $1,
                             updated_at = NOW()
WHERE id = $2 AND user_id = $3;

-- name: ClubParticipantCheckInviteApproval :one
SELECT cp.id FROM club_participants as cp
INNER JOIN clubs as c ON c.id = cp.club_id
WHERE cp.id = $1 AND c.user_id = $2 AND user_approval = TRUE AND club_approval IS NULL LIMIT 1;

-- name: ClubParticipantCheckJoinApproval :one
SELECT id FROM club_participants WHERE id = $1 AND user_id = $2 AND user_approval IS NULL AND club_approval = TRUE LIMIT 1;

-- name: ClubParticipantCheckOne :one
SELECT id FROM club_participants  WHERE user_approval = TRUE AND club_approval = TRUE AND club_id = $1 AND user_id = $2;

-- name: ClubParticipantKick :exec
UPDATE club_participants SET club_approval = FALSE,
                             updated_at = NOW()
WHERE id = $1 AND club_id = $2;

-- name: ClubParticipantCheckByOwner :one
SELECT cp.id FROM club_participants as cp
INNER JOIN clubs as c ON c.id = cp.club_id
WHERE cp.user_approval = TRUE AND cp.club_approval = TRUE AND cp.club_id = $1 AND c.user_id = $2 AND cp.id = $3;