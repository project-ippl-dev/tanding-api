-- name: ClubCreate :one
INSERT INTO clubs(name, logo, phone, short_name, user_id, updated_at) VALUES ($1, $2, $3, $4, $5, NOW()) RETURNING id;
