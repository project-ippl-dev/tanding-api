package bracket

import (
	"context"
	"database/sql"

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

const orderBracketFetchByClassEventID = `-- name: OrderBracketFetchByClassEventID :many
SELECT ob.id, ob.rank, ob.order_by,
       array(SELECT u.name FROM event_participants as ep
           INNER JOIN users as u ON u.id = ep.user_id
           WHERE ep.event_registration_id = ob.event_registration_id
        ) as participants, c.name as club_name, ob.event_registration_id, c.logo as club_logo
FROM order_brackets as ob
INNER JOIN clubs as c ON c.id = ob.club_id
WHERE ob.class_event_id = $1 ORDER BY ob.order_by
`

type OrderBracketFetchByClassEventIDRow struct {
	ID                  uuid.UUID `json:"id"`
	Rank                int16     `json:"rank"`
	OrderBy             int16     `json:"order_by"`
	Participants        []string  `json:"participants"`
	ClubName            string    `json:"club_name"`
	ClubLogo            string    `json:"club_logo"`
	EventRegistrationID uuid.UUID `json:"event_registration_id"`
}

func (r *RawRepository) OrderBracketFetchByClassEventID(ctx context.Context, classEventID uuid.UUID) ([]OrderBracketFetchByClassEventIDRow, error) {
	rows, err := r.db.QueryContext(ctx, orderBracketFetchByClassEventID, classEventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []OrderBracketFetchByClassEventIDRow{}
	for rows.Next() {
		var i OrderBracketFetchByClassEventIDRow
		if err := rows.Scan(
			&i.ID,
			&i.Rank,
			&i.OrderBy,
			pq.Array(&i.Participants),
			&i.ClubName,
			&i.EventRegistrationID,
			&i.ClubLogo,
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

const bracketParticipantFetchByEventBracketID = `-- name: BracketParticipantFetchByEventBracketID :many
SELECT bp.id, er.club_id, c.name as club_name, bp.type, array(
        SELECT u.name FROM event_participants as ep
                               INNER JOIN users as u ON u.id = ep.user_id
        WHERE ep.event_registration_id = er.id
    ) as participants, bp.is_bye, bp.event_registration_id
FROM bracket_participants as bp
LEFT JOIN event_registrations as er ON er.id = bp.event_registration_id
LEFT JOIN clubs as c ON c.id  = er.club_id
WHERE bp.event_bracket_id = $1
ORDER BY CASE bp.type
             WHEN 'home' THEN 1
             WHEN 'away' THEN 2
             END
`

type bracketParticipantFetchByEventBracketIDRow struct {
	ID                  int64              `json:"id"`
	ClubID              uuid.UUID          `json:"club_id"`
	ClubName            sql.NullString     `json:"club_name"`
	Type                db.ParticipantType `json:"type"`
	Participants        []string           `json:"participants"`
	IsBye               bool               `json:"is_bye"`
	EventRegistrationID uuid.UUID          `json:"event_registration_id"`
}

func (r *RawRepository) BracketParticipantFetchByEventBracketID(ctx context.Context, eventBracketID uuid.UUID) ([]bracketParticipantFetchByEventBracketIDRow, error) {
	rows, err := r.db.QueryContext(ctx, bracketParticipantFetchByEventBracketID, eventBracketID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []bracketParticipantFetchByEventBracketIDRow{}
	for rows.Next() {
		var i bracketParticipantFetchByEventBracketIDRow
		if err := rows.Scan(
			&i.ID,
			&i.ClubID,
			&i.ClubName,
			&i.Type,
			pq.Array(&i.Participants),
			&i.IsBye,
			&i.EventRegistrationID,
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

const rankFetchByClassEventID = `-- name: RankFetchByClassEventID :many
SELECT r.id, r.rank, r.point, c.name as club_name ,array(
    SELECT u.name FROM event_participants as ep
    INNER JOIN users as u ON u.id = ep.user_id
    WHERE ep.event_registration_id = er.id
) as participants, c.logo
FROM ranks as r
INNER JOIN event_registrations as er ON er.id = r.event_registration_id
INNER JOIN clubs as c ON c.id = er.club_id
WHERE r.class_event_id = $1
ORDER BY r.rank
`

type RankFetchByClassEventIDRow struct {
	ID           uuid.UUID `json:"id"`
	Rank         int16     `json:"rank"`
	Point        int32     `json:"point"`
	ClubName     string    `json:"club_name"`
	ClubLogo     string    `json:"club_logo"`
	Participants []string  `json:"participants"`
}

func (r *RawRepository) RankFetchByClassEventID(ctx context.Context, classEventID uuid.UUID) ([]RankFetchByClassEventIDRow, error) {
	rows, err := r.db.QueryContext(ctx, rankFetchByClassEventID, classEventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []RankFetchByClassEventIDRow{}
	for rows.Next() {
		var i RankFetchByClassEventIDRow
		if err := rows.Scan(
			&i.ID,
			&i.Rank,
			&i.Point,
			&i.ClubName,
			pq.Array(&i.Participants),
			&i.ClubLogo,
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
