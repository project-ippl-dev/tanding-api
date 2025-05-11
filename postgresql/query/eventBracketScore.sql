-- name: EventScoreCheckOneByBracketID :one
SELECT id FROM event_scores WHERE event_bracket_id = $1;

-- name: EventScoreCreate :exec
INSERT INTO event_scores(event_bracket_id, home_round1, home_round2, home_round3, home_extra, home_total,
                         away_round1, away_round2, away_round3, away_extra, away_total, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW());

-- name: EventScoreUpdate :exec
UPDATE event_scores SET home_round1 = $1,
                        home_round2 = $2,
                        home_round3 = $3,
                        home_extra = $4,
                        home_total = $5,
                        away_round1 = $6,
                        away_round2 = $7,
                        away_round3 = $8,
                        away_extra = $9,
                        away_total = $10,
                        updated_at = NOW()
WHERE event_bracket_id = $11;

-- name: EventScoreFetchOneByBracketID :one
SELECT id, home_round1, home_round2, home_round3, home_extra, home_total,
       away_round1, away_round2, away_round3, away_extra, away_total
FROM event_scores
WHERE event_bracket_id = $1;

-- name: EventScoreFetchHomeByBracketID :one
SELECT id, home_round1, home_round2, home_round3, home_extra, home_total
FROM event_scores
WHERE event_bracket_id = $1;

-- name: EventScoreFetchAwayByBracketID :one
SELECT id, away_round1, away_round2, away_round3, away_extra, away_total
FROM event_scores
WHERE event_bracket_id = $1;
