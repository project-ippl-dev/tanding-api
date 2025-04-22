-- name: DocumentCreate :exec
INSERT INTO documents(user_id, birth_certificate, family_card, user_identity, belt_certificate, elementary_certificate,
                      middle_certificate, high_certificate, bachelor_certificate, master_certificate, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7 ,$8 ,$9 ,$10, NOW());

-- name: DocumentFetchAll :many
SELECT id, user_id, birth_certificate, family_card, user_identity, belt_certificate, elementary_certificate,
       middle_certificate, high_certificate, bachelor_certificate, master_certificate FROM documents
WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: DocumentCountAll :one
SELECT COUNT(id) FROM documents WHERE user_id = $1 LIMIT 1;

-- name: DocumentCheckOne :one
SELECT id FROM documents WHERE id = $1 AND user_id = $2 LIMIT 1;

-- name: DocumentUpdate :exec
UPDATE documents SET birth_certificate = $1,
                     family_card = $2,
                     user_identity = $3,
                     belt_certificate = $4,
                     elementary_certificate = $5,
                     middle_certificate = $6,
                     high_certificate = $7,
                     bachelor_certificate = $8,
                     master_certificate = $9,
                     updated_at = NOW()
WHERE id = $10 ;

-- name: DocumentDelete :exec
DELETE FROM documents WHERE id = $1;