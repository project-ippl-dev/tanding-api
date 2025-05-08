-- name: ClassCreate :exec
INSERT INTO classes(sport_id, name, class_competition_rule_id, type, match_type, updated_at)
VALUES($1, $2, $3, $4, $5, NOW());

-- name: ClassFetchAll :many
SELECT c.id, c.name, c.class_competition_rule_id, ccr.name as class_rule_name, c.type, c.sport_id, c.match_type, s.name as sport_name FROM classes as c
INNER JOIN sports as s ON s.id = c.sport_id
INNER JOIN class_competition_rules as ccr ON ccr.id = c.class_competition_rule_id
ORDER BY c.created_at DESC LIMIT $1 OFFSET $2;

-- name: ClassCountAll :one
SELECT COUNT(id) FROM classes LIMIT 1;

-- name: ClassFetchBySportID :many
SELECT c.id, c.name, c.class_competition_rule_id, ccr.name as class_rule_name, c.type, c.sport_id, c.match_type, s.name as sport_name FROM classes as c
INNER JOIN sports as s ON s.id = c.sport_id
INNER JOIN class_competition_rules as ccr ON ccr.id = c.class_competition_rule_id
WHERE c.sport_id = $1 ORDER BY c.created_at DESC LIMIT $2 OFFSET $3;

-- name: ClassCountBySportID :one
SELECT COUNT(id) FROM classes WHERE sport_id = $1 LIMIT 1;

-- name: ClassCheckOne :one
SELECT id, name, type, match_type FROM classes WHERE id = $1 LIMIT 1;

-- name: ClassUpdate :exec
UPDATE classes SET sport_id = $1,
                   name = $2,
                   class_competition_rule_id = $3,
                   type = $4,
                   match_type = $5,
                   updated_at = NOW()
WHERE id = $6;

-- name: ClassDelete :exec
DELETE FROM classes WHERE id = $1;