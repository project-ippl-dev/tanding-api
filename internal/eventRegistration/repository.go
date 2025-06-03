package eventRegistration

import (
	"context"
	"database/sql"
	"strings"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/project-ippl-dev/tanding-api/internal/db"
)

type RawRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) RawRepository {
	return RawRepository{db: db}
}

const eventRegistrationFetchAll = `-- name: EventRegistrationFetchAll :many
SELECT er.id, c.name as class_name,cl.name as club_name, er.status, er.club_id, er.class_event_id, ce.price FROM event_registrations as er
INNER JOIN class_events as ce ON ce.id = er.class_event_id
INNER JOIN classes as c ON c.id = ce.class_id
INNER JOIN clubs as cl ON cl.id = er.club_id
WHERE er.event_id = $1
ORDER BY er.created_at DESC LIMIT $2 OFFSET $3
`

const eventRegistrationFetchByClubID = `-- name: EventRegistrationFetchByClubID :many
SELECT er.id, c.name as class_name,cl.name as club_name, er.status, er.club_id, er.class_event_id, ce.price FROM event_registrations as er
INNER JOIN class_events as ce ON ce.id = er.class_event_id
INNER JOIN classes as c ON c.id = ce.class_id
INNER JOIN clubs as cl ON cl.id = er.club_id
WHERE er.event_id = $1 AND er.club_id = $2
ORDER BY er.created_at DESC LIMIT $3 OFFSET $4
`

const eventRegistrationFetchByClassEventID = `-- name: EventRegistrationFetchByClassEventID :many
SELECT er.id, c.name as class_name,cl.name as club_name, er.status, er.club_id, er.class_event_id, ce.price FROM event_registrations as er
INNER JOIN class_events as ce ON ce.id = er.class_event_id
INNER JOIN classes as c ON c.id = ce.class_id
INNER JOIN clubs as cl ON cl.id = er.club_id
WHERE er.event_id = $1 AND er.class_event_id = $2
ORDER BY er.created_at DESC LIMIT $3 OFFSET $4
`

const eventRegistrationFetchByStatus = `-- name: EventRegistrationFetchByStatus :many
SELECT er.id, c.name as class_name,cl.name as club_name, er.status, er.club_id, er.class_event_id, ce.price FROM event_registrations as er
INNER JOIN class_events as ce ON ce.id = er.class_event_id
INNER JOIN classes as c ON c.id = ce.class_id
INNER JOIN clubs as cl ON cl.id = er.club_id
WHERE er.event_id = $1 AND er.status = $2
ORDER BY er.created_at DESC LIMIT $3 OFFSET $4
`

const eventRegistrationFetchByUserID = `-- name: EventRegistrationFetchByUserID :many
SELECT er.id, c.name as class_name,cl.name as club_name, er.status, er.club_id, er.class_event_id, ce.price FROM event_registrations as er
INNER JOIN class_events as ce ON ce.id = er.class_event_id
INNER JOIN classes as c ON c.id = ce.class_id
INNER JOIN clubs as cl ON cl.id = er.club_id
WHERE er.event_id = $1 AND cl.user_id = $2
ORDER BY er.created_at DESC LIMIT $3 OFFSET $4
`

const eventRegistrationFetchByClubIDAndClassEventID = `-- name: EventRegistrationFetchByClubIDAndClassEventID :many
SELECT er.id, c.name as class_name,cl.name as club_name, er.status, er.club_id, er.class_event_id, ce.price FROM event_registrations as er
INNER JOIN class_events as ce ON ce.id = er.class_event_id
INNER JOIN classes as c ON c.id = ce.class_id
INNER JOIN clubs as cl ON cl.id = er.club_id
WHERE er.event_id = $1 AND er.club_id = $2 AND er.class_event_id = $3
ORDER BY er.created_at DESC LIMIT $4 OFFSET $5
`

const eventRegistrationFetchByClubIDAndStatus = `-- name: EventRegistrationFetchByClubIDAndStatus :many
SELECT er.id, c.name as class_name,cl.name as club_name, er.status, er.club_id, er.class_event_id, ce.price FROM event_registrations as er
INNER JOIN class_events as ce ON ce.id = er.class_event_id
INNER JOIN classes as c ON c.id = ce.class_id
INNER JOIN clubs as cl ON cl.id = er.club_id
WHERE er.event_id = $1 AND er.club_id = $2 AND er.status = $3
ORDER BY er.created_at DESC LIMIT $4 OFFSET $5
`

const eventRegistrationFetchByClubIDAndUserID = `-- name: EventRegistrationFetchByClubIDAndUserID :many
SELECT er.id, c.name as class_name,cl.name as club_name, er.status, er.club_id, er.class_event_id, ce.price FROM event_registrations as er
INNER JOIN class_events as ce ON ce.id = er.class_event_id
INNER JOIN classes as c ON c.id = ce.class_id
INNER JOIN clubs as cl ON cl.id = er.club_id
WHERE er.event_id = $1 AND er.club_id = $2 AND cl.user_id = $3
ORDER BY er.created_at DESC LIMIT $4 OFFSET $5
`

const eventRegistrationFetchByClassEventIDAndStatus = `-- name: EventRegistrationFetchByClassEventIDAndStatus :many
SELECT er.id, c.name as class_name,cl.name as club_name, er.status, er.club_id, er.class_event_id, ce.price FROM event_registrations as er
INNER JOIN class_events as ce ON ce.id = er.class_event_id
INNER JOIN classes as c ON c.id = ce.class_id
INNER JOIN clubs as cl ON cl.id = er.club_id
WHERE er.event_id = $1 AND er.class_event_id = $2 AND er.status = $3
ORDER BY er.created_at DESC LIMIT $4 OFFSET $5
`

const eventRegistrationFetchByClassEventIDAndUserID = `-- name: EventRegistrationFetchByClassEventIDAndUserID :many
SELECT er.id, c.name as class_name,cl.name as club_name, er.status, er.club_id, er.class_event_id, ce.price FROM event_registrations as er
INNER JOIN class_events as ce ON ce.id = er.class_event_id
INNER JOIN classes as c ON c.id = ce.class_id
INNER JOIN clubs as cl ON cl.id = er.club_id
WHERE er.event_id = $1 AND er.class_event_id = $2 AND er.user_id = $3
ORDER BY er.created_at DESC LIMIT $4 OFFSET $5
`

const eventRegistrationFetchByStatusAndUserID = `-- name: EventRegistrationFetchByStatusAndUserID :many
SELECT er.id, c.name as class_name,cl.name as club_name, er.status, er.club_id, er.class_event_id, ce.price FROM event_registrations as er
INNER JOIN class_events as ce ON ce.id = er.class_event_id
INNER JOIN classes as c ON c.id = ce.class_id
INNER JOIN clubs as cl ON cl.id = er.club_id
WHERE er.event_id = $1 AND er.status = $2 AND er.user_id = $3
ORDER BY er.created_at DESC LIMIT $4 OFFSET $5
`

const eventRegistrationFetchByClubIDClassEventIDAndStatus = `-- name: EventRegistrationFetchByClubIDClassEventIDAndStatus :many
SELECT er.id, c.name as class_name,cl.name as club_name, er.status, er.club_id, er.class_event_id, ce.price FROM event_registrations as er
INNER JOIN class_events as ce ON ce.id = er.class_event_id
INNER JOIN classes as c ON c.id = ce.class_id
INNER JOIN clubs as cl ON cl.id = er.club_id
WHERE er.event_id = $1 AND er.club_id = $2 AND er.class_event_id = $3 AND er.status = $4
ORDER BY er.created_at DESC LIMIT $5 OFFSET $6
`

const eventRegistrationFetchByClubIDClassEventIDAndUserID = `-- name: EventRegistrationFetchByClubIDClassEventIDAndUserID :many
SELECT er.id, c.name as class_name,cl.name as club_name, er.status, er.club_id, er.class_event_id, ce.price FROM event_registrations as er
INNER JOIN class_events as ce ON ce.id = er.class_event_id
INNER JOIN classes as c ON c.id = ce.class_id
INNER JOIN clubs as cl ON cl.id = er.club_id
WHERE er.event_id = $1 AND er.club_id = $2 AND er.class_event_id = $3 AND cl.user_id = $4
ORDER BY er.created_at DESC LIMIT $5 OFFSET $6
`

const eventRegistrationFetchByClubIDStatusAndUserID = `-- name: EventRegistrationFetchByClubIDStatusAndUserID :many
SELECT er.id, c.name as class_name,cl.name as club_name, er.status, er.club_id, er.class_event_id, ce.price FROM event_registrations as er
INNER JOIN class_events as ce ON ce.id = er.class_event_id
INNER JOIN classes as c ON c.id = ce.class_id
INNER JOIN clubs as cl ON cl.id = er.club_id
WHERE er.event_id = $1 AND er.club_id = $2 AND er.status = $3 AND cl.user_id = $4
ORDER BY er.created_at DESC LIMIT $5 OFFSET $6
`

const eventRegistrationFetchByClassEventIDStatusAndUserID = `-- name: EventRegistrationFetchByClassEventIDStatusAndUserID :many
SELECT er.id, c.name as class_name,cl.name as club_name, er.status, er.club_id, er.class_event_id, ce.price FROM event_registrations as er
INNER JOIN class_events as ce ON ce.id = er.class_event_id
INNER JOIN classes as c ON c.id = ce.class_id
INNER JOIN clubs as cl ON cl.id = er.club_id
WHERE er.event_id = $1 AND er.class_event_id = $2 AND er.status = $3 AND cl.user_id = $4
ORDER BY er.created_at DESC LIMIT $5 OFFSET $6
`

const eventRegistrationFetchByFilterAll = `-- name: EventRegistrationFetchByClassEventIDAndStatus :many
SELECT er.id, c.name as class_name,cl.name as club_name, er.status, er.club_id, er.class_event_id, ce.price FROM event_registrations as er
INNER JOIN class_events as ce ON ce.id = er.class_event_id
INNER JOIN classes as c ON c.id = ce.class_id
INNER JOIN clubs as cl ON cl.id = er.club_id
WHERE er.event_id = $1 AND er.club_id = $2 AND er.class_event_id = $3 AND er.status = $4 AND cl.user_id = $5
ORDER BY er.created_at DESC LIMIT $6 OFFSET $7
`

type fetchQueryParams struct {
	FetchAllParams
}

type fetchAllDBParams struct {
	fetchQueryParams
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type FetchAllRow struct {
	ID           uuid.UUID                  `json:"id"`
	ClassEventID uuid.UUID                  `json:"class_event_id"`
	ClassName    string                     `json:"class_name"`
	ClubID       uuid.UUID                  `json:"club_id"`
	ClubName     string                     `json:"club_name"`
	Status       db.EventRegistrationStatus `json:"status"`
	Price        int32                      `json:"price"`
}

func (r *RawRepository) setQueryStatus(arg fetchQueryParams) string {
	var data []string
	if arg.ClubID != "" {
		data = append(data, "club_id")
	}
	if arg.ClassEventID != "" {
		data = append(data, "class_event_id")
	}
	if arg.Status != "" {
		data = append(data, "status")
	}
	if arg.UserID != "" {
		data = append(data, "user_id")
	}

	return strings.Join(data, "-")
}

func (r *RawRepository) EventRegistrationFetchAll(ctx context.Context, arg fetchAllDBParams) ([]FetchAllRow, error) {
	var rows *sql.Rows
	status := r.setQueryStatus(arg.fetchQueryParams)
	switch status {
	default:
		result, err := r.db.QueryContext(ctx, eventRegistrationFetchAll,
			arg.EventID,
			arg.Limit, arg.Offset)
		if err != nil {
			return nil, err
		}
		rows = result
	case "club_id":
		result, err := r.db.QueryContext(ctx, eventRegistrationFetchByClubID,
			arg.EventID, arg.ClubID,
			arg.Limit, arg.Offset)
		if err != nil {
			return nil, err
		}
		rows = result
	case "class_event_id":
		result, err := r.db.QueryContext(ctx, eventRegistrationFetchByClassEventID,
			arg.EventID, arg.ClassEventID,
			arg.Limit, arg.Offset)
		if err != nil {
			return nil, err
		}
		rows = result
	case "status":
		result, err := r.db.QueryContext(ctx, eventRegistrationFetchByStatus,
			arg.EventID, arg.Status,
			arg.Limit, arg.Offset)
		if err != nil {
			return nil, err
		}
		rows = result
	case "user_id":
		result, err := r.db.QueryContext(ctx, eventRegistrationFetchByUserID,
			arg.EventID, arg.UserID,
			arg.Limit, arg.Offset)
		if err != nil {
			return nil, err
		}
		rows = result
	case "club_id-class_event_id":
		result, err := r.db.QueryContext(ctx, eventRegistrationFetchByClubIDAndClassEventID,
			arg.EventID, arg.ClubID, arg.ClassEventID,
			arg.Limit, arg.Offset)
		if err != nil {
			return nil, err
		}
		rows = result
	case "club_id-status":
		result, err := r.db.QueryContext(ctx, eventRegistrationFetchByClubIDAndStatus,
			arg.EventID, arg.ClubID, arg.Status,
			arg.Limit, arg.Offset)
		if err != nil {
			return nil, err
		}
		rows = result
	case "club_id-user_id":
		result, err := r.db.QueryContext(ctx, eventRegistrationFetchByClubIDAndUserID,
			arg.EventID, arg.ClubID, arg.UserID,
			arg.Limit, arg.Offset)
		if err != nil {
			return nil, err
		}
		rows = result
	case "class_event_id-status":
		result, err := r.db.QueryContext(ctx, eventRegistrationFetchByClassEventIDAndStatus,
			arg.EventID, arg.ClassEventID, arg.Status,
			arg.Limit, arg.Offset)
		if err != nil {
			return nil, err
		}
		rows = result
	case "class_event_id-user_id":
		result, err := r.db.QueryContext(ctx, eventRegistrationFetchByClassEventIDAndUserID,
			arg.EventID, arg.ClassEventID, arg.UserID,
			arg.Limit, arg.Offset)
		if err != nil {
			return nil, err
		}
		rows = result
	case "status-user_id":
		result, err := r.db.QueryContext(ctx, eventRegistrationFetchByStatusAndUserID,
			arg.EventID, arg.Status, arg.UserID,
			arg.Limit, arg.Offset)
		if err != nil {
			return nil, err
		}
		rows = result
	case "club_id-class_event_id-status":
		result, err := r.db.QueryContext(ctx, eventRegistrationFetchByClubIDClassEventIDAndStatus,
			arg.EventID, arg.ClubID, arg.ClassEventID, arg.Status,
			arg.Limit, arg.Offset)
		if err != nil {
			return nil, err
		}
		rows = result
	case "club_id-class_event_id-user_id":
		result, err := r.db.QueryContext(ctx, eventRegistrationFetchByClubIDClassEventIDAndUserID,
			arg.EventID, arg.ClubID, arg.ClassEventID, arg.UserID,
			arg.Limit, arg.Offset)
		if err != nil {
			return nil, err
		}
		rows = result
	case "club_id-status-user_id":
		result, err := r.db.QueryContext(ctx, eventRegistrationFetchByClubIDStatusAndUserID,
			arg.EventID, arg.ClubID, arg.Status, arg.UserID,
			arg.Limit, arg.Offset)
		if err != nil {
			return nil, err
		}
		rows = result
	case "class_event_id-status-user_id":
		result, err := r.db.QueryContext(ctx, eventRegistrationFetchByClassEventIDStatusAndUserID,
			arg.EventID, arg.ClassEventID, arg.Status, arg.UserID,
			arg.Limit, arg.Offset)
		if err != nil {
			return nil, err
		}
		rows = result
	case "club_id-class_event_id-status-user_id":
		result, err := r.db.QueryContext(ctx, eventRegistrationFetchByFilterAll,
			arg.EventID, arg.ClubID, arg.ClassEventID, arg.Status, arg.UserID,
			arg.Limit, arg.Offset)
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
			&i.ClassName,
			&i.ClubName,
			&i.Status,
			&i.ClubID,
			&i.ClassEventID,
			&i.Price,
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

const eventRegistrationCountAll = `-- name: EventRegistrationCountAll :one
SELECT COUNT(er.id) FROM event_registrations as er
INNER JOIN class_events as ce ON ce.id = er.class_event_id
INNER JOIN classes as c ON c.id = ce.class_id
INNER JOIN clubs as cl ON cl.id = er.club_id
WHERE er.event_id = $1
`

const eventRegistrationCountByClubID = `-- name: EventRegistrationCountByClubID :one
SELECT COUNT(er.id) FROM event_registrations as er
INNER JOIN class_events as ce ON ce.id = er.class_event_id
INNER JOIN classes as c ON c.id = ce.class_id
INNER JOIN clubs as cl ON cl.id = er.club_id
WHERE er.event_id = $1 AND er.club_id = $2
`

const eventRegistrationCountByClassEventID = `-- name: EventRegistrationCountByClassEventID :one
SELECT COUNT(er.id) FROM event_registrations as er
INNER JOIN class_events as ce ON ce.id = er.class_event_id
INNER JOIN classes as c ON c.id = ce.class_id
INNER JOIN clubs as cl ON cl.id = er.club_id
WHERE er.event_id = $1 AND er.class_event_id = $2
`

const eventRegistrationCountByStatus = `-- name: EventRegistrationCountByStatus :one
SELECT COUNT(er.id) FROM event_registrations as er
INNER JOIN class_events as ce ON ce.id = er.class_event_id
INNER JOIN classes as c ON c.id = ce.class_id
INNER JOIN clubs as cl ON cl.id = er.club_id
WHERE er.event_id = $1 AND er.status = $2
`

const eventRegistrationCountByUserID = `-- name: EventRegistrationCountByUserID :one
SELECT COUNT(er.id) FROM event_registrations as er
INNER JOIN class_events as ce ON ce.id = er.class_event_id
INNER JOIN classes as c ON c.id = ce.class_id
INNER JOIN clubs as cl ON cl.id = er.club_id
WHERE er.event_id = $1 AND cl.user_id = $2
`

const eventRegistrationCountByClubIDAndClassEventID = `-- name: EventRegistrationCountByClubIDAndClassEventID :one
SELECT COUNT(er.id) FROM event_registrations as er
INNER JOIN class_events as ce ON ce.id = er.class_event_id
INNER JOIN classes as c ON c.id = ce.class_id
INNER JOIN clubs as cl ON cl.id = er.club_id
WHERE er.event_id = $1 AND er.club_id = $2 AND er.class_event_id = $3
`

const eventRegistrationCountByClubIDAndStatus = `-- name: EventRegistrationCountByClubIDAndStatus :one
SELECT COUNT(er.id) FROM event_registrations as er
INNER JOIN class_events as ce ON ce.id = er.class_event_id
INNER JOIN classes as c ON c.id = ce.class_id
INNER JOIN clubs as cl ON cl.id = er.club_id
WHERE er.event_id = $1 AND er.club_id = $2 AND er.status = $3
`

const eventRegistrationCountByClubIDAndUserID = `-- name: EventRegistrationCountByClubIDAndUserID :one
SELECT COUNT(er.id) FROM event_registrations as er
INNER JOIN class_events as ce ON ce.id = er.class_event_id
INNER JOIN classes as c ON c.id = ce.class_id
INNER JOIN clubs as cl ON cl.id = er.club_id
WHERE er.event_id = $1 AND er.club_id = $2 AND cl.user_id = $3
`

const eventRegistrationCountByClassEventIDAndStatus = `-- name: EventRegistrationCountByClassEventIDAndStatus :one
SELECT COUNT(er.id) FROM event_registrations as er
INNER JOIN class_events as ce ON ce.id = er.class_event_id
INNER JOIN classes as c ON c.id = ce.class_id
INNER JOIN clubs as cl ON cl.id = er.club_id
WHERE er.event_id = $1 AND er.class_event_id = $2 AND er.status = $3
`

const eventRegistrationCountByClassEventIDAndUserID = `-- name: EventRegistrationCountByClassEventIDAndUserID :one
SELECT COUNT(er.id) FROM event_registrations as er
INNER JOIN class_events as ce ON ce.id = er.class_event_id
INNER JOIN classes as c ON c.id = ce.class_id
INNER JOIN clubs as cl ON cl.id = er.club_id
WHERE er.event_id = $1 AND er.class_event_id = $2 AND cl.user_id = $3
`

const eventRegistrationCountByStatusAndUserID = `-- name: EventRegistrationCountByStatusAndUserID :one
SELECT COUNT(er.id) FROM event_registrations as er
INNER JOIN class_events as ce ON ce.id = er.class_event_id
INNER JOIN classes as c ON c.id = ce.class_id
INNER JOIN clubs as cl ON cl.id = er.club_id
WHERE er.event_id = $1 AND er.status = $2 AND cl.user_id = $3
`

const eventRegistrationCountByClubIDClassEventIDAndStatus = `-- name: EventRegistrationCountByClubIDClassEventIDAndStatus :one
SELECT COUNT(er.id) FROM event_registrations as er
INNER JOIN class_events as ce ON ce.id = er.class_event_id
INNER JOIN classes as c ON c.id = ce.class_id
INNER JOIN clubs as cl ON cl.id = er.club_id
WHERE er.event_id = $1 AND er.club_id = $2 AND er.class_event_id = $3 AND er.status = $4
`

const eventRegistrationCountByClubIDClassEventIDAndUserID = `-- name: EventRegistrationCountByClubIDClassEventIDAndUserID :one
SELECT COUNT(er.id) FROM event_registrations as er
INNER JOIN class_events as ce ON ce.id = er.class_event_id
INNER JOIN classes as c ON c.id = ce.class_id
INNER JOIN clubs as cl ON cl.id = er.club_id
WHERE er.event_id = $1 AND er.club_id = $2 AND er.class_event_id = $3 AND cl.user_id = $4
`

const eventRegistrationCountByClubIDStatusAndUserID = `-- name: EventRegistrationCountByClubIDStatusAndUserID :one
SELECT COUNT(er.id) FROM event_registrations as er
INNER JOIN class_events as ce ON ce.id = er.class_event_id
INNER JOIN classes as c ON c.id = ce.class_id
INNER JOIN clubs as cl ON cl.id = er.club_id
WHERE er.event_id = $1 AND er.club_id = $2 AND er.status = $3 AND cl.user_id = $4
`

const eventRegistrationCountByClassEventIDStatusAndUserID = `-- name: EventRegistrationCountByClassEventIDStatusAndUserID :one
SELECT COUNT(er.id) FROM event_registrations as er
INNER JOIN class_events as ce ON ce.id = er.class_event_id
INNER JOIN classes as c ON c.id = ce.class_id
INNER JOIN clubs as cl ON cl.id = er.club_id
WHERE er.event_id = $1 AND er.class_event_id = $2 AND er.status = $3 AND cl.user_id = $4
`

const eventRegistrationCountByFilterAll = `-- name: EventRegistrationCountByFilterAll :one
SELECT COUNT(er.id) FROM event_registrations as er
INNER JOIN class_events as ce ON ce.id = er.class_event_id
INNER JOIN classes as c ON c.id = ce.class_id
INNER JOIN clubs as cl ON cl.id = er.club_id
WHERE er.event_id = $1 AND er.club_id = $2 AND er.class_event_id = $3 AND er.status = $4 AND cl.user_id = $5
`

func (r *RawRepository) EventRegistrationCountAll(ctx context.Context, arg fetchQueryParams) (int64, error) {
	var row *sql.Row
	status := r.setQueryStatus(arg)
	switch status {
	default:
		row = r.db.QueryRowContext(ctx, eventRegistrationCountAll, arg.EventID)
	case "club_id":
		row = r.db.QueryRowContext(ctx, eventRegistrationCountByClubID, arg.EventID, arg.ClubID)
	case "class_event_id":
		row = r.db.QueryRowContext(ctx, eventRegistrationCountByClassEventID, arg.EventID, arg.ClassEventID)
	case "status":
		row = r.db.QueryRowContext(ctx, eventRegistrationCountByStatus, arg.EventID, arg.Status)
	case "user_id":
		row = r.db.QueryRowContext(ctx, eventRegistrationCountByUserID, arg.EventID, arg.UserID)
	case "club_id-class_event_id":
		row = r.db.QueryRowContext(ctx, eventRegistrationCountByClubIDAndClassEventID, arg.EventID, arg.ClubID, arg.ClassEventID)
	case "club_id-status":
		row = r.db.QueryRowContext(ctx, eventRegistrationCountByClubIDAndStatus, arg.EventID, arg.ClubID, arg.Status)
	case "club_id-user_id":
		row = r.db.QueryRowContext(ctx, eventRegistrationCountByClubIDAndUserID, arg.EventID, arg.ClubID, arg.UserID)
	case "class_event_id-status":
		row = r.db.QueryRowContext(ctx, eventRegistrationCountByClassEventIDAndStatus, arg.EventID, arg.ClassEventID, arg.Status)
	case "class_event_id-user_id":
		row = r.db.QueryRowContext(ctx, eventRegistrationCountByClassEventIDAndUserID, arg.EventID, arg.ClassEventID, arg.UserID)
	case "status-user_id":
		row = r.db.QueryRowContext(ctx, eventRegistrationCountByStatusAndUserID, arg.EventID, arg.Status, arg.UserID)
	case "club_id-class_event_id-status":
		row = r.db.QueryRowContext(ctx, eventRegistrationCountByClubIDClassEventIDAndStatus, arg.EventID, arg.ClubID, arg.ClassEventID, arg.Status)
	case "club_id-class_event_id-user_id":
		row = r.db.QueryRowContext(ctx, eventRegistrationCountByClubIDClassEventIDAndUserID, arg.EventID, arg.ClubID, arg.ClassEventID, arg.UserID)
	case "club_id-status-user_id":
		row = r.db.QueryRowContext(ctx, eventRegistrationCountByClubIDStatusAndUserID, arg.EventID, arg.ClubID, arg.Status, arg.UserID)
	case "class_event_id-status-user_id":
		row = r.db.QueryRowContext(ctx, eventRegistrationCountByClassEventIDStatusAndUserID, arg.EventID, arg.ClassEventID, arg.Status, arg.UserID)
	case "club_id-class_event_id-status-user_id":
		row = r.db.QueryRowContext(ctx, eventRegistrationCountByFilterAll, arg.EventID, arg.ClubID, arg.ClassEventID, arg.Status, arg.UserID)
	}

	var count int64
	err := row.Scan(&count)
	return count, err
}

const userFetchInID = `-- name: UserFetchInID :many
SELECT id, name, gender, can_participate FROM users WHERE id = ANY($1)
`

type UserFetchInIDRow struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Gender         string    `json:"gender"`
	CanParticipate bool      `json:"can_participate"`
}

func (r *RawRepository) UserFetchInID(ctx context.Context, id []string) ([]UserFetchInIDRow, error) {
	rows, err := r.db.QueryContext(ctx, userFetchInID, pq.Array(id))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []UserFetchInIDRow{}
	for rows.Next() {
		var i UserFetchInIDRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Gender,
			&i.CanParticipate,
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
