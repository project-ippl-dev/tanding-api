-- name: AccomplishmentCreate :exec
INSERT INTO accomplishments(user_id, title, level, ranking, category, sport, description, file_url, month, year, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW());

-- name: AccomplishmentFetchAll :many
SELECT id, user_id, title, level, ranking, category, sport, description, file_url, month, year FROM accomplishments
WHERE user_id = $1 ORDER BY year DESC, month DESC LIMIT $2 OFFSET $3;

-- name: AccomplishmentCountAll :one
SELECT COUNT(id) FROM accomplishments WHERE user_id = $1 LIMIT 1;

-- name: AccomplishmentCheckOne :one
SELECT id FROM accomplishments WHERE id = $1 AND user_id = $2 LIMIT 1;

-- name: AccomplishmentUpdate :exec
UPDATE accomplishments SET title = $1,
                           level = $2,
                           ranking = $3,
                           category = $4,
                           sport = $5,
                           description = $6,
                           file_url = $7,
                           month = $8,
                           year = $9,
                           updated_at = NOW()
WHERE id = $10;

-- name: AccomplishmentDelete :exec
DELETE FROM accomplishments WHERE id = $1;