package rank

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type RawRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) RawRepository {
	return RawRepository{db: db}
}

const eventParticipantFetchExcludeRegistrationID = `-- name: EventParticipantFetchExcludeRegistrationID :many
SELECT ep.id, ep.event_registration_id, ep.user_id, u.name FROM event_participants as ep
INNER JOIN event_registrations as er ON er.id = ep.event_registration_id
INNER JOIN users as u ON u.id = ep.user_id
WHERE er.class_event_id = $1 AND NOT(ep.event_registration_id = ANY($2)) AND er.status = 'approved'
`

type EventParticipantFetchExcludeRegistrationIDParams struct {
	EventRegistrationID []uuid.UUID
	ClassEventID        uuid.UUID
}

type EventParticipantFetchExcludeRegistrationIDRow struct {
	ID                  int64     `json:"id"`
	EventRegistrationID uuid.UUID `json:"event_registration_id"`
	UserID              uuid.UUID `json:"user_id"`
	Name                string    `json:"name"`
}

func (r *RawRepository) EventParticipantFetchExcludeRegistrationID(ctx context.Context, arg EventParticipantFetchExcludeRegistrationIDParams) ([]EventParticipantFetchExcludeRegistrationIDRow, error) {
	rows, err := r.db.QueryContext(ctx, eventParticipantFetchExcludeRegistrationID, arg.ClassEventID, pq.Array(arg.EventRegistrationID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []EventParticipantFetchExcludeRegistrationIDRow{}
	for rows.Next() {
		var i EventParticipantFetchExcludeRegistrationIDRow
		if err := rows.Scan(
			&i.ID,
			&i.EventRegistrationID,
			&i.UserID,
			&i.Name,
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

const rankFetchPowerList = `-- name: RankFetchPowerList :many
SELECT c.id, c.name, (COALESCE(SUM(r.point),0)) as total_point, (COALESCE(COUNT(er.id),0)) AS total_participate,
       c.logo
FROM ranks as r
         INNER JOIN event_registrations as er ON er.id = r.event_registration_id
         INNER JOIN clubs as c ON c.id = er.club_id
GROUP BY c.id, c.name
ORDER BY total_point DESC, total_participate DESC
LIMIT $1 OFFSET $2
`

const rankFetchPowerListBySportID = `-- name: RankFetchPowerList :many
SELECT c.id, c.name, (COALESCE(SUM(r.point),0)) as total_point, (COALESCE(COUNT(er.id),0)) AS total_participate,
       c.logo as photo
FROM ranks as r
         INNER JOIN event_registrations as er ON er.id = r.event_registration_id
         INNER JOIN clubs as c ON c.id = er.club_id
GROUP BY c.id, c.name, r.sport_id
HAVING r.sport_id = $1
ORDER BY total_point DESC, total_participate DESC
LIMIT $2 OFFSET $3
`

type fetchPowerListParams struct {
	SportID string
	Limit   int32 `json:"limit"`
	Offset  int32 `json:"offset"`
}

type fetchPowerListRow struct {
	ID               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	TotalPoint       int64     `json:"total_point"`
	TotalParticipate int64     `json:"total_participate"`
	Photo            string    `json:"photo"`
}

func (r *RawRepository) RankFetchPowerList(ctx context.Context, arg fetchPowerListParams) ([]fetchPowerListRow, error) {
	var rows *sql.Rows
	if arg.SportID == "" {
		result, err := r.db.QueryContext(ctx, rankFetchPowerList, arg.Limit, arg.Offset)
		if err != nil {
			return nil, err
		}
		rows = result
	} else {
		result, err := r.db.QueryContext(ctx, rankFetchPowerListBySportID, arg.SportID, arg.Limit, arg.Offset)
		if err != nil {
			return nil, err
		}
		rows = result
	}
	defer rows.Close()
	items := []fetchPowerListRow{}
	for rows.Next() {
		var i fetchPowerListRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.TotalPoint,
			&i.TotalParticipate,
			&i.Photo,
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

const rankCountPowerList = `-- name: RankCountPowerList :many
SELECT COUNT(*) FROM (SELECT c.id, c.name, (COALESCE(SUM(r.point),0)) as total_point, (COALESCE(COUNT(er.id),0)) AS total_participate
FROM ranks as r
         INNER JOIN event_registrations as er ON er.id = r.event_registration_id
         INNER JOIN clubs as c ON c.id = er.club_id
GROUP BY c.id, c.name) as result
`

const rankCountPowerListBySportID = `-- name: RankCountPowerList :many
SELECT COUNT(*) FROM (SELECT c.id, c.name, (COALESCE(SUM(r.point),0)) as total_point, (COALESCE(COUNT(er.id),0)) AS total_participate
FROM ranks as r
         INNER JOIN event_registrations as er ON er.id = r.event_registration_id
         INNER JOIN clubs as c ON c.id = er.club_id
GROUP BY c.id, c.name, r.sport_id
HAVING r.sport_id = $1) as result
`

func (r *RawRepository) RankCountPowerList(ctx context.Context, sportID string) (int64, error) {
	var row *sql.Row
	if sportID == "" {
		row = r.db.QueryRowContext(ctx, rankCountPowerList)
	} else {
		row = r.db.QueryRowContext(ctx, rankCountPowerListBySportID, sportID)
	}
	var count int64
	err := row.Scan(&count)
	return count, err
}

const rankFetchAllPointUser = `-- name: RankFetchAllPointUser :many
SELECT u.id, u.name, COALESCE(SUM(r.point),0) as total_point, r.club_id, c.name as club_name, 
       (COALESCE(COUNT(er.id),0)) AS total_participate, u.photo
FROM ranks as r
INNER JOIN event_registrations as er ON er.id = r.event_registration_id
INNER JOIN event_participants as ep ON ep.event_registration_id = er.id
INNER JOIN clubs as c ON c.id = r.club_id
INNER JOIN users as u ON u.id = ep.user_id
GROUP BY u.id, r.club_id, c.name
ORDER BY total_point DESC
LIMIT $1 OFFSET $2
`

const rankFetchAllPointUserBySportID = `-- name: RankFetchAllPointUser :many
SELECT u.id, u.name, COALESCE(SUM(r.point),0) as total_point, r.club_id, c.name as club_name, 
       (COALESCE(COUNT(er.id),0)) AS total_participate, u.photo
FROM ranks as r
INNER JOIN event_registrations as er ON er.id = r.event_registration_id
INNER JOIN event_participants as ep ON ep.event_registration_id = er.id
INNER JOIN clubs as c ON c.id = r.club_id
INNER JOIN users as u ON u.id = ep.user_id
GROUP BY u.id, r.club_id, c.name, r.sport_id
HAVING r.sport_id = $1
ORDER BY total_point DESC
LIMIT $2 OFFSET $3
`

type fetchAllPointUserParams struct {
	SportID string `json:"sport_id"`
	Limit   int32  `json:"limit"`
	Offset  int32  `json:"offset"`
}

type fetchAllPointUserRow struct {
	ID               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	ClubID           uuid.UUID `json:"club_id"`
	ClubName         string    `json:"club_name"`
	TotalPoint       int64     `json:"total_point"`
	TotalParticipate int64     `json:"total_participate"`
	Photo            string    `json:"photo"`
}

func (r *RawRepository) RankFetchAllPointUser(ctx context.Context, arg fetchAllPointUserParams) ([]fetchAllPointUserRow, error) {
	var rows *sql.Rows
	if arg.SportID == "" {
		result, err := r.db.QueryContext(ctx, rankFetchAllPointUser, arg.Limit, arg.Offset)
		if err != nil {
			return nil, err
		}
		rows = result
	} else {
		result, err := r.db.QueryContext(ctx, rankFetchAllPointUserBySportID, arg.SportID, arg.Limit, arg.Offset)
		if err != nil {
			return nil, err
		}
		rows = result
	}
	defer rows.Close()
	items := []fetchAllPointUserRow{}
	for rows.Next() {
		var i fetchAllPointUserRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.TotalPoint,
			&i.ClubID,
			&i.ClubName,
			&i.TotalParticipate,
			&i.Photo,
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

const rankCountAllPointUser = `-- name: RankCountAllPointUser :one
SELECT COUNT(*) FROM (SELECT u.id, u.name, COALESCE(SUM(r.point),0) as point FROM ranks as r
INNER JOIN event_registrations as er ON er.id = r.event_registration_id
INNER JOIN event_participants as ep ON ep.event_registration_id = er.id
INNER JOIN users as u ON u.id = ep.user_id
GROUP BY u.id, r.club_id
ORDER BY point DESC) as total
`

const rankCountAllPointUserBySportID = `-- name: RankCountAllPointUser :one
SELECT COUNT(*) FROM (SELECT u.id, u.name, COALESCE(SUM(r.point),0) as total_point, r.club_id, c.name as club_name, 
       (COALESCE(COUNT(er.id),0)) AS total_participate
FROM ranks as r
INNER JOIN event_registrations as er ON er.id = r.event_registration_id
INNER JOIN event_participants as ep ON ep.event_registration_id = er.id
INNER JOIN clubs as c ON c.id = r.club_id
INNER JOIN users as u ON u.id = ep.user_id
GROUP BY u.id, r.club_id, c.name, r.sport_id
HAVING r.sport_id = $1
ORDER BY total_point DESC) as total
`

func (r *RawRepository) RankCountAllPointUser(ctx context.Context, sportID string) (int64, error) {
	var row *sql.Row
	if sportID == "" {
		row = r.db.QueryRowContext(ctx, rankCountAllPointUser)
	} else {
		row = r.db.QueryRowContext(ctx, rankCountAllPointUserBySportID, sportID)
	}
	var count int64
	err := row.Scan(&count)
	return count, err
}
