package club

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"strings"
)

type RawRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) RawRepository {
	return RawRepository{db: db}
}

const clubFetchAll = `-- name: ClubFetchAll :many
SELECT clubs.id, clubs.name, clubs.logo, users.name as owner FROM clubs
INNER JOIN users ON users.id = clubs.user_id
ORDER BY clubs.created_at DESC LIMIT $1 OFFSET $2
`

const clubFetchBySportID = `-- name: ClubFetchBySportID :many
SELECT clubs.id, clubs.name, clubs.logo, users.name as owner FROM clubs
INNER JOIN club_sport as cp ON cp.club_id = clubs.id
INNER JOIN users ON users.id = clubs.user_id
WHERE cp.sport_id = $1 ORDER BY clubs.created_at DESC
LIMIT $2 OFFSET $3
`

type fetchAllParams struct {
	Limit   int32  `json:"limit"`
	Offset  int32  `json:"offset"`
	SportID string `json:"sport_id"`
}

type FetchAllRow struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Logo  string    `json:"logo"`
	Owner string    `json:"owner"`
}

func (r *RawRepository) ClubFetchAll(ctx context.Context, arg fetchAllParams) ([]FetchAllRow, error) {
	var rows *sql.Rows
	if arg.SportID == "" {
		result, err := r.db.QueryContext(ctx, clubFetchAll, arg.Limit, arg.Offset)
		if err != nil {
			return nil, err
		}
		rows = result
	} else {
		result, err := r.db.QueryContext(ctx, clubFetchBySportID, arg.SportID, arg.Limit, arg.Offset)
		if err != nil {
			return nil, err
		}
		rows = result
	}
	defer rows.Close()
	items := []FetchAllRow{}
	for rows.Next() {
		var i FetchAllRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Logo,
			&i.Owner,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const clubCountAll = `-- name: ClubCountAll :one
SELECT COUNT(id) FROM clubs
`

const clubCountBySportID = `-- name: ClubCountBySportID :one
SELECT COUNT(clubs.id) FROM clubs
INNER JOIN club_sport as cp ON cp.club_id = clubs.id
WHERE cp.sport_id = $1
`

func (r *RawRepository) ClubCountAll(ctx context.Context, sportID string) (int64, error) {
	var row *sql.Row
	if sportID == "" {
		row = r.db.QueryRowContext(ctx, clubCountAll)
	} else {
		row = r.db.QueryRowContext(ctx, clubCountBySportID, sportID)
	}

	var count int64
	err := row.Scan(&count)
	return count, err
}

const clubParticipantFetchAll = `-- name: ClubParticipantFetchAll :many
SELECT cp.id, u.name, u.can_participate, u.id as user_id FROM club_participants as cp
INNER JOIN users as u ON cp.user_id = u.id
WHERE cp.club_id = $1 AND cp.user_approval = TRUE AND cp.club_approval = TRUE
ORDER BY cp.id DESC LIMIT $2 OFFSET $3
`

const clubParticipantFetchBySportID = `-- name: ClubParticipantFetchBySportID :many
SELECT cp.id, u.name, u.can_participate, u.id as user_id FROM club_participants as cp
INNER JOIN users as u ON cp.user_id = u.id
WHERE cp.club_id = $1 AND cp.sport_id = $2 AND cp.user_approval = TRUE AND cp.club_approval = TRUE
ORDER BY cp.id DESC LIMIT $3 OFFSET $4
`

const clubParticipantFetchByParticipate = `-- name: ClubParticipantFetchByParticipate :many
SELECT cp.id, u.name, u.can_participate, u.id as user_id FROM club_participants as cp
INNER JOIN users as u ON cp.user_id = u.id
WHERE cp.club_id = $1 AND u.can_participate = $2 AND cp.user_approval = TRUE AND cp.club_approval = TRUE
ORDER BY cp.id DESC LIMIT $3 OFFSET $4
`

const clubParticipantFetchByFilterAll = `-- name: ClubParticipantFetchByFilterAll :many
SELECT cp.id, u.name, u.can_participate, u.id as user_id FROM club_participants as cp
INNER JOIN users as u ON cp.user_id = u.id
WHERE cp.club_id = $1 AND cp.sport_id = $2 AND u.can_participate = $3 AND cp.user_approval = TRUE 
AND cp.club_approval = TRUE
ORDER BY cp.id DESC LIMIT $4 OFFSET $5
`

type participantQueryParams struct {
	SportID        string    `json:"sport_id"`
	CanParticipate *bool     `json:"can_participate"`
	ClubID         uuid.UUID `json:"club_id"`
}

type participantFetchAllParams struct {
	participantQueryParams
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type ParticipantFetchAllRow struct {
	ID             int64     `json:"id"`
	UserID         uuid.UUID `json:"user_id"`
	Name           string    `json:"name"`
	CanParticipate bool      `json:"can_participate"`
}

func (r *RawRepository) setQueryStatus(arg participantQueryParams) string {
	var result []string
	if arg.SportID != "" {
		result = append(result, "sport")
	}
	if arg.CanParticipate != nil {
		result = append(result, "participate")
	}

	return strings.Join(result, "-")
}

func (r *RawRepository) ClubParticipantFetchAll(ctx context.Context, arg participantFetchAllParams) ([]ParticipantFetchAllRow, error) {
	var rows *sql.Rows
	status := r.setQueryStatus(arg.participantQueryParams)
	switch status {
	default:
		result, err := r.db.QueryContext(ctx, clubParticipantFetchAll, arg.ClubID, arg.Limit, arg.Offset)
		if err != nil {
			return nil, err
		}
		rows = result
	case "sport":
		result, err := r.db.QueryContext(ctx, clubParticipantFetchBySportID, arg.ClubID, arg.SportID, arg.Limit, arg.Offset)
		if err != nil {
			return nil, err
		}
		rows = result
	case "participate":
		result, err := r.db.QueryContext(ctx, clubParticipantFetchByParticipate, arg.ClubID, arg.CanParticipate, arg.Limit, arg.Offset)
		if err != nil {
			return nil, err
		}
		rows = result
	case "sport-participate":
		result, err := r.db.QueryContext(ctx, clubParticipantFetchByFilterAll, arg.ClubID, arg.SportID, arg.CanParticipate, arg.Limit, arg.Offset)
		if err != nil {
			return nil, err
		}
		rows = result
	}
	defer rows.Close()
	items := []ParticipantFetchAllRow{}
	for rows.Next() {
		var i ParticipantFetchAllRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.CanParticipate,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const clubParticipantCountAll = `-- name: ClubParticipantCountAll :one
SELECT COUNT(id) FROM club_participants WHERE club_id = $1 
                                          AND user_approval = TRUE AND club_approval = TRUE
`

const clubParticipantCountBySportID = `-- name: ClubParticipantCountBySportID :one
SELECT COUNT(id) FROM club_participants WHERE club_id = $1 AND sport_id = $2
                                          AND user_approval = TRUE AND club_approval = TRUE
`

const clubParticipantCountByCanParticipate = `-- name: ClubParticipantCountBySportID :one
SELECT COUNT(cp.id) FROM club_participants as cp 
INNER JOIN users as u ON cp.user_id = u.id
WHERE cp.club_id = $1 AND u.can_participate = $2
AND user_approval = TRUE AND club_approval = TRUE
`

const clubParticipantCountByFilterAll = `-- name: ClubParticipantCountByFilterAll :one
SELECT COUNT(cp.id) FROM club_participants as cp 
INNER JOIN users as u ON cp.user_id = u.id
WHERE cp.club_id = $1 AND cp.sport_id = $2 AND u.can_participate = $3
AND user_approval = TRUE AND club_approval = TRUE
`

func (r *RawRepository) ClubParticipantCountAll(ctx context.Context, arg participantQueryParams) (int64, error) {
	var row *sql.Row
	status := r.setQueryStatus(arg)
	switch status {
	default:
		row = r.db.QueryRowContext(ctx, clubParticipantCountAll, arg.ClubID)
	case "sport":
		row = r.db.QueryRowContext(ctx, clubParticipantCountBySportID, arg.ClubID, arg.SportID)
	case "participate":
		row = r.db.QueryRowContext(ctx, clubParticipantCountByCanParticipate, arg.ClubID, arg.CanParticipate)
	case "sport-participate":
		row = r.db.QueryRowContext(ctx, clubParticipantCountByFilterAll, arg.ClubID, arg.SportID, arg.CanParticipate)
	}
	var count int64
	err := row.Scan(&count)
	return count, err
}
