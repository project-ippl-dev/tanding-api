-- name: ClassRuleCreate :exec
INSERT INTO class_competition_rules(name, male, female, total, updated_at) VALUES ($1, $2, $3, $4, NOW());

-- name: ClassRuleFetchAll :many
SELECT id, name, male, female, total FROM class_competition_rules ORDER BY name LIMIT $1 OFFSET $2;

-- name: ClassRuleCountAll :one
SELECT COUNT(id) FROM class_competition_rules;

-- name: ClassRuleFetchOne :one
SELECT id, name, male, female, total FROM class_competition_rules WHERE id = $1 LIMIT 1;

-- name: ClassRuleUpdate :exec
UPDATE class_competition_rules SET name = $1,
                                   male = $2,
                                   female = $3,
                                   total = $4,
                                   updated_at = NOW()
WHERE id = $5;

-- name: ClassRuleDelete :exec
DELETE FROM class_competition_rules WHERE id = $1;