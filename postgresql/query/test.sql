-- name: TestFetchOne :one
SELECT * FROM test WHERE id = $1 LIMIT 1;