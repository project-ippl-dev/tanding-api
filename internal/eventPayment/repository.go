package eventPayment

import (
	"context"
	"database/sql"
	"strings"
	"time"

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

const eventPaymentFetchAll = `-- name: EventPaymentFetchAll :many
SELECT DISTINCT ep.id, ep.unique_number, ep.payment_link, ep.status, c.name as club_name, ep.created_at, 
	er.club_id, ep.total, o.name as club_owner, u.name as user_name, a.name as admin_name, e.name as event_name, e.id as event_id
FROM event_payment_receipts as ep
INNER JOIN event_registrations as er ON er.event_payment_receipt_id = ep.id
INNER JOIN events as e ON e.id = er.event_id
INNER JOIN clubs as c ON c.id = er.club_id 
INNER JOIN users as u ON u.id = ep.user_id
LEFT JOIN users as a ON a.id = ep.admin_id
INNER JOIN users as o ON o.id = c.user_id
ORDER BY ep.created_at DESC LIMIT $1 OFFSET $2
`

const eventPaymentFetchByEventID = `-- name: EventPaymentFetchByEventID :many
SELECT DISTINCT ep.id, ep.unique_number, ep.payment_link, ep.status, c.name as club_name, ep.created_at, 
	er.club_id, ep.total, o.name as club_owner, u.name as user_name, a.name as admin_name, e.name as event_name, e.id as event_id
FROM event_payment_receipts as ep
INNER JOIN event_registrations as er ON er.event_payment_receipt_id = ep.id
INNER JOIN events as e ON e.id = er.event_id
INNER JOIN clubs as c ON c.id = er.club_id 
INNER JOIN users as u ON u.id = ep.user_id
LEFT JOIN users as a ON a.id = ep.admin_id
INNER JOIN users as o ON o.id = c.user_id
WHERE er.event_id = $1 
ORDER BY ep.created_at DESC LIMIT $2 OFFSET $3
`

const eventPaymentFetchByStatus = `-- name: EventPaymentFetchByStatus :many
SELECT DISTINCT ep.id, ep.unique_number, ep.payment_link, ep.status, c.name as club_name, ep.created_at, 
	er.club_id, ep.total, o.name as club_owner, u.name as user_name, a.name as admin_name, e.name as event_name, e.id as event_id
FROM event_payment_receipts as ep
INNER JOIN event_registrations as er ON er.event_payment_receipt_id = ep.id
INNER JOIN events as e ON e.id = er.event_id
INNER JOIN clubs as c ON c.id = er.club_id 
INNER JOIN users as u ON u.id = ep.user_id
LEFT JOIN users as a ON a.id = ep.admin_id
INNER JOIN users as o ON o.id = c.user_id
WHERE ep.status = $1 
ORDER BY ep.created_at DESC LIMIT $2 OFFSET $3
`

const eventPaymentFetchByClubID = `-- name: EventPaymentFetchByClubID :many
SELECT DISTINCT ep.id, ep.unique_number, ep.payment_link, ep.status, c.name as club_name, ep.created_at, 
	er.club_id, ep.total, o.name as club_owner, u.name as user_name, a.name as admin_name, e.name as event_name, e.id as event_id
FROM event_payment_receipts as ep
INNER JOIN event_registrations as er ON er.event_payment_receipt_id = ep.id
INNER JOIN events as e ON e.id = er.event_id
INNER JOIN clubs as c ON c.id = er.club_id 
INNER JOIN users as u ON u.id = ep.user_id
LEFT JOIN users as a ON a.id = ep.admin_id
INNER JOIN users as o ON o.id = c.user_id
WHERE er.club_id = any($1) 
ORDER BY ep.created_at DESC LIMIT $2 OFFSET $3
`

const eventPaymentFetchByDate = `-- name: EventPaymentFetchByDate :many
SELECT DISTINCT ep.id, ep.unique_number, ep.payment_link, ep.status, c.name as club_name, ep.created_at, 
	er.club_id, ep.total, o.name as club_owner, u.name as user_name, a.name as admin_name, e.name as event_name, e.id as event_id
FROM event_payment_receipts as ep
INNER JOIN event_registrations as er ON er.event_payment_receipt_id = ep.id
INNER JOIN events as e ON e.id = er.event_id
INNER JOIN clubs as c ON c.id = er.club_id 
INNER JOIN users as u ON u.id = ep.user_id
LEFT JOIN users as a ON a.id = ep.admin_id
INNER JOIN users as o ON o.id = c.user_id
WHERE ep.created_at BETWEEN DATE($1) AND DATE($2) 
ORDER BY ep.created_at DESC LIMIT $3 OFFSET $4
`

const eventPaymentFetchByEventIDAndStatus = `-- name: EventPaymentFetchByEventIDAndStatus :many
SELECT DISTINCT ep.id, ep.unique_number, ep.payment_link, ep.status, c.name as club_name, ep.created_at, 
	er.club_id, ep.total, o.name as club_owner, u.name as user_name, a.name as admin_name, e.name as event_name, e.id as event_id
FROM event_payment_receipts as ep
INNER JOIN event_registrations as er ON er.event_payment_receipt_id = ep.id
INNER JOIN events as e ON e.id = er.event_id
INNER JOIN clubs as c ON c.id = er.club_id 
INNER JOIN users as u ON u.id = ep.user_id
LEFT JOIN users as a ON a.id = ep.admin_id
INNER JOIN users as o ON o.id = c.user_id
WHERE er.event_id = $1 AND ep.status = $2 
ORDER BY ep.created_at DESC LIMIT $3 OFFSET $4
`

const eventPaymentFetchByEventIDAndClubID = `-- name: EventPaymentFetchByEventIDAndStatus :many
SELECT DISTINCT ep.id, ep.unique_number, ep.payment_link, ep.status, c.name as club_name, ep.created_at, 
	er.club_id, ep.total, o.name as club_owner, u.name as user_name, a.name as admin_name, e.name as event_name, e.id as event_id
FROM event_payment_receipts as ep
INNER JOIN event_registrations as er ON er.event_payment_receipt_id = ep.id
INNER JOIN events as e ON e.id = er.event_id
INNER JOIN clubs as c ON c.id = er.club_id 
INNER JOIN users as u ON u.id = ep.user_id
LEFT JOIN users as a ON a.id = ep.admin_id
INNER JOIN users as o ON o.id = c.user_id
WHERE er.event_id = $1 AND er.club_id = any($2) 
ORDER BY ep.created_at DESC LIMIT $3 OFFSET $4
`

const eventPaymentFetchByEventIDAndDate = `-- name: EventPaymentFetchByEventIDAndDate :many
SELECT DISTINCT ep.id, ep.unique_number, ep.payment_link, ep.status, c.name as club_name, ep.created_at, 
	er.club_id, ep.total, o.name as club_owner, u.name as user_name, a.name as admin_name, e.name as event_name, e.id as event_id
FROM event_payment_receipts as ep
INNER JOIN event_registrations as er ON er.event_payment_receipt_id = ep.id
INNER JOIN events as e ON e.id = er.event_id
INNER JOIN clubs as c ON c.id = er.club_id 
INNER JOIN users as u ON u.id = ep.user_id
LEFT JOIN users as a ON a.id = ep.admin_id
INNER JOIN users as o ON o.id = c.user_id
WHERE er.event_id = $1 AND ep.created_at BETWEEN DATE($2) AND DATE($3) 
ORDER BY ep.created_at DESC LIMIT $4 OFFSET $5
`

const eventPaymentFetchByStatusAndClubID = `-- name: EventPaymentFetchByStatusAndClubID :many
SELECT DISTINCT ep.id, ep.unique_number, ep.payment_link, ep.status, c.name as club_name, ep.created_at, 
	er.club_id, ep.total, o.name as club_owner, u.name as user_name, a.name as admin_name, e.name as event_name, e.id as event_id
FROM event_payment_receipts as ep
INNER JOIN event_registrations as er ON er.event_payment_receipt_id = ep.id
INNER JOIN events as e ON e.id = er.event_id
INNER JOIN clubs as c ON c.id = er.club_id 
INNER JOIN users as u ON u.id = ep.user_id
LEFT JOIN users as a ON a.id = ep.admin_id
INNER JOIN users as o ON o.id = c.user_id
WHERE ep.status = $1 AND er.club_id = any($2) 
ORDER BY ep.created_at DESC LIMIT $3 OFFSET $4
`

const eventPaymentFetchByStatusAndDate = `-- name: EventPaymentFetchByStatusAndDate :many
SELECT DISTINCT ep.id, ep.unique_number, ep.payment_link, ep.status, c.name as club_name, ep.created_at, 
	er.club_id, ep.total, o.name as club_owner, u.name as user_name, a.name as admin_name, e.name as event_name, e.id as event_id
FROM event_payment_receipts as ep
INNER JOIN event_registrations as er ON er.event_payment_receipt_id = ep.id
INNER JOIN events as e ON e.id = er.event_id
INNER JOIN clubs as c ON c.id = er.club_id 
INNER JOIN users as u ON u.id = ep.user_id
LEFT JOIN users as a ON a.id = ep.admin_id
INNER JOIN users as o ON o.id = c.user_id
WHERE ep.status = $1 AND created_at BETWEEN DATE($2) AND DATE($3) 
ORDER BY ep.created_at DESC LIMIT $4 OFFSET $5
`

const eventPaymentFetchByClubIDAndDate = `-- name: EventPaymentFetchByClubIDAndDate :many
SELECT DISTINCT ep.id, ep.unique_number, ep.payment_link, ep.status, c.name as club_name, ep.created_at, 
	er.club_id, ep.total, o.name as club_owner, u.name as user_name, a.name as admin_name, e.name as event_name, e.id as event_id
FROM event_payment_receipts as ep
INNER JOIN event_registrations as er ON er.event_payment_receipt_id = ep.id
INNER JOIN events as e ON e.id = er.event_id
INNER JOIN clubs as c ON c.id = er.club_id 
INNER JOIN users as u ON u.id = ep.user_id
LEFT JOIN users as a ON a.id = ep.admin_id
INNER JOIN users as o ON o.id = c.user_id
WHERE er.club_id = any($1) AND ep.created_at BETWEEN DATE($2) AND DATE($3)
ORDER BY ep.created_at DESC LIMIT $4 OFFSET $5
`

const eventPaymentFetchByEventIDStatusAndClubID = `-- name: EventPaymentFetchByEventIDStatusAndClubID :many
SELECT DISTINCT ep.id, ep.unique_number, ep.payment_link, ep.status, c.name as club_name, ep.created_at, 
	er.club_id, ep.total, o.name as club_owner, u.name as user_name, a.name as admin_name, e.name as event_name, e.id as event_id
FROM event_payment_receipts as ep
INNER JOIN event_registrations as er ON er.event_payment_receipt_id = ep.id
INNER JOIN events as e ON e.id = er.event_id
INNER JOIN clubs as c ON c.id = er.club_id 
INNER JOIN users as u ON u.id = ep.user_id
LEFT JOIN users as a ON a.id = ep.admin_id
INNER JOIN users as o ON o.id = c.user_id
WHERE er.event_id = $1 AND ep.status = $2 AND er.club_id = any($3) 
ORDER BY ep.created_at DESC LIMIT $4 OFFSET $5
`

const eventPaymentFetchByEventIDStatusAndDate = `-- name: EventPaymentFetchByEventIDStatusAndDate :many
SELECT DISTINCT ep.id, ep.unique_number, ep.payment_link, ep.status, c.name as club_name, ep.created_at, 
	er.club_id, ep.total, o.name as club_owner, u.name as user_name, a.name as admin_name, e.name as event_name, e.id as event_id
FROM event_payment_receipts as ep
INNER JOIN event_registrations as er ON er.event_payment_receipt_id = ep.id
INNER JOIN events as e ON e.id = er.event_id
INNER JOIN clubs as c ON c.id = er.club_id 
INNER JOIN users as u ON u.id = ep.user_id
LEFT JOIN users as a ON a.id = ep.admin_id
INNER JOIN users as o ON o.id = c.user_id
WHERE er.event_id = $1 AND ep.status = $2 AND BETWEEN DATE($3) AND DATE($4)
ORDER BY ep.created_at DESC LIMIT $5 OFFSET $6
`

const eventPaymentFetchByEventIDClubIDAndDate = `-- name: EventPaymentFetchByEventIDClubIDAndDate :many
SELECT DISTINCT ep.id, ep.unique_number, ep.payment_link, ep.status, c.name as club_name, ep.created_at, 
	er.club_id, ep.total, o.name as club_owner, u.name as user_name, a.name as admin_name, e.name as event_name, e.id as event_id
FROM event_payment_receipts as ep
INNER JOIN event_registrations as er ON er.event_payment_receipt_id = ep.id
INNER JOIN events as e ON e.id = er.event_id
INNER JOIN clubs as c ON c.id = er.club_id 
INNER JOIN users as u ON u.id = ep.user_id
LEFT JOIN users as a ON a.id = ep.admin_id
INNER JOIN users as o ON o.id = c.user_id
WHERE er.event_id = $1 AND er.club_id = any($2) AND ep.created_at BETWEEN DATE($3) AND DATE($4)
ORDER BY ep.created_at DESC LIMIT $5 OFFSET $6
`

const eventPaymentFetchByStatusClubIDAndDate = `-- name: EventPaymentFetchByStatusClubIDAndDate :many
SELECT DISTINCT ep.id, ep.unique_number, ep.payment_link, ep.status, c.name as club_name, ep.created_at, 
	er.club_id, ep.total, o.name as club_owner, u.name as user_name, a.name as admin_name, e.name as event_name, e.id as event_id
FROM event_payment_receipts as ep
INNER JOIN event_registrations as er ON er.event_payment_receipt_id = ep.id
INNER JOIN events as e ON e.id = er.event_id
INNER JOIN clubs as c ON c.id = er.club_id 
INNER JOIN users as u ON u.id = ep.user_id
LEFT JOIN users as a ON a.id = ep.admin_id
INNER JOIN users as o ON o.id = c.user_id
WHERE er.club_id = any($1) AND ep.created_at BETWEEN DATE($2) AND DATE($3) AND ep.status = $4 
ORDER BY ep.created_at DESC LIMIT $5 OFFSET $6
`

const eventPaymentFetchByEventIDFilterAll = `-- name: EventPaymentFetchByEventIDFilterAll :many
SELECT DISTINCT ep.id, ep.unique_number, ep.payment_link, ep.status, c.name as club_name, ep.created_at, 
	er.club_id, ep.total, o.name as club_owner, u.name as user_name, a.name as admin_name, e.name as event_name, e.id as event_id
FROM event_payment_receipts as ep
INNER JOIN event_registrations as er ON er.event_payment_receipt_id = ep.id
INNER JOIN events as e ON e.id = er.event_id
INNER JOIN clubs as c ON c.id = er.club_id 
INNER JOIN users as u ON u.id = ep.user_id
LEFT JOIN users as a ON a.id = ep.admin_id
INNER JOIN users as o ON o.id = c.user_id
WHERE er.event_id = $1 AND er.club_id = any($2) AND ep.status = $3 AND ep.created_at BETWEEN DATE($4) AND DATE($5)
ORDER BY ep.created_at DESC LIMIT $6 OFFSET $7
`

type queryPaymentParams struct {
	EventID string                `json:"event_id"`
	Status  db.EventReceiptStatus `json:"status"`
	ClubID  []string              `json:"club_id"`
	Start   time.Time             `json:"start"`
	End     time.Time             `json:"end"`
}

type fetchByEventIDParams struct {
	queryPaymentParams
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type fetchAllRow struct {
	ID           uuid.UUID             `json:"id"`
	EventID      uuid.UUID             `json:"event_id"`
	EventName    string                `json:"event_name"`
	UniqueNumber int16                 `json:"unique_number"`
	PaymentLink  string                `json:"payment_link"`
	Status       db.EventReceiptStatus `json:"status"`
	ClubID       uuid.UUID             `json:"club_id"`
	ClubName     string                `json:"club_name"`
	UserName     string                `json:"user_name"`
	AdminName    sql.NullString        `json:"admin_name"`
	ClubOwner    string                `json:"club_owner"`
	Total        int32                 `json:"total"`
	CreatedAt    time.Time             `json:"created_at"`
}

func (r *RawRepository) setQueryStatus(arg queryPaymentParams) string {
	var result []string
	if arg.EventID != "" {
		result = append(result, "event")
	}
	if arg.Status != "" {
		result = append(result, "status")
	}
	if arg.ClubID != nil {
		result = append(result, "club_id")
	}
	if !arg.Start.IsZero() && !arg.End.IsZero() {
		result = append(result, "date")
	}

	return strings.Join(result, "-")
}

func (r *RawRepository) EventPaymentFetchAll(ctx context.Context, arg fetchByEventIDParams) ([]fetchAllRow, error) {
	var rows *sql.Rows
	status := r.setQueryStatus(arg.queryPaymentParams)
	switch status {
	default:
		result, err := r.db.QueryContext(ctx, eventPaymentFetchAll, arg.Limit, arg.Offset)
		if err != nil {
			return nil, err
		}
		rows = result
	case "event":
		result, err := r.db.QueryContext(ctx, eventPaymentFetchByEventID, arg.EventID, arg.Limit, arg.Offset)
		if err != nil {
			return nil, err
		}
		rows = result
	case "status":
		result, err := r.db.QueryContext(ctx, eventPaymentFetchByStatus, arg.Status, arg.Limit, arg.Offset)
		if err != nil {
			return nil, err
		}
		rows = result
	case "club_id":
		result, err := r.db.QueryContext(ctx, eventPaymentFetchByClubID, pq.Array(arg.ClubID), arg.Limit, arg.Offset)
		if err != nil {
			return nil, err
		}
		rows = result
	case "date":
		result, err := r.db.QueryContext(ctx, eventPaymentFetchByDate, arg.Start, arg.End, arg.Limit, arg.Offset)
		if err != nil {
			return nil, err
		}
		rows = result
	case "event-status":
		result, err := r.db.QueryContext(ctx, eventPaymentFetchByEventIDAndStatus, arg.EventID, arg.Status, arg.Limit, arg.Offset)
		if err != nil {
			return nil, err
		}
		rows = result
	case "event-club":
		result, err := r.db.QueryContext(ctx, eventPaymentFetchByEventIDAndClubID, arg.EventID, pq.Array(arg.ClubID), arg.Limit, arg.Offset)
		if err != nil {
			return nil, err
		}
		rows = result
	case "event-date":
		result, err := r.db.QueryContext(ctx, eventPaymentFetchByEventIDAndDate, arg.EventID, arg.Start, arg.End, arg.Limit, arg.Offset)
		if err != nil {
			return nil, err
		}
		rows = result
	case "status-club_id":
		result, err := r.db.QueryContext(ctx, eventPaymentFetchByStatusAndClubID, arg.Status, pq.Array(arg.ClubID), arg.Limit, arg.Offset)
		if err != nil {
			return nil, err
		}
		rows = result
	case "status-date":
		result, err := r.db.QueryContext(ctx, eventPaymentFetchByStatusAndDate, arg.Status, arg.Start, arg.End, arg.Limit, arg.Offset)
		if err != nil {
			return nil, err
		}
		rows = result
	case "club_id-date":
		result, err := r.db.QueryContext(ctx, eventPaymentFetchByClubIDAndDate, pq.Array(arg.ClubID), arg.Start, arg.End, arg.Limit, arg.Offset)
		if err != nil {
			return nil, err
		}
		rows = result
	case "event-status-club":
		result, err := r.db.QueryContext(ctx, eventPaymentFetchByEventIDStatusAndClubID, arg.EventID, arg.Status, pq.Array(arg.ClubID), arg.Limit, arg.Offset)
		if err != nil {
			return nil, err
		}
		rows = result
	case "event-status-date":
		result, err := r.db.QueryContext(ctx, eventPaymentFetchByEventIDStatusAndDate, arg.EventID, arg.Status, arg.Start, arg.End, arg.Limit, arg.Offset)
		if err != nil {
			return nil, err
		}
		rows = result
	case "event-club-date":
		result, err := r.db.QueryContext(ctx, eventPaymentFetchByEventIDClubIDAndDate, arg.EventID, arg.ClubID, arg.Start, arg.End, arg.Limit, arg.Offset)
		if err != nil {
			return nil, err
		}
		rows = result
	case "status-club-date":
		result, err := r.db.QueryContext(ctx, eventPaymentFetchByStatusClubIDAndDate, arg.Status, pq.Array(arg.ClubID), arg.Start, arg.End, arg.Limit, arg.Offset)
		if err != nil {
			return nil, err
		}
		rows = result
	case "event-status-club_id-date":
		result, err := r.db.QueryContext(ctx, eventPaymentFetchByEventIDFilterAll, arg.EventID, arg.Status, pq.Array(arg.ClubID), arg.Start, arg.End, arg.Limit, arg.Offset)
		if err != nil {
			return nil, err
		}
		rows = result
	}

	defer rows.Close()
	items := []fetchAllRow{}
	for rows.Next() {
		var i fetchAllRow
		if err := rows.Scan(
			&i.ID,
			&i.UniqueNumber,
			&i.PaymentLink,
			&i.Status,
			&i.ClubName,
			&i.CreatedAt,
			&i.ClubID,
			&i.Total,
			&i.ClubOwner,
			&i.UserName,
			&i.AdminName,
			&i.EventName,
			&i.EventID,
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

const eventPaymentCountByAll = `-- name: EventPaymentCountByAll :one
SELECT COUNT(*) FROM (SELECT DISTINCT ep.id FROM event_payment_receipts as ep
INNER JOIN event_registrations as er ON er.event_payment_receipt_id = ep.id
INNER JOIN clubs as c ON c.id = er.club_id) AS count
`

const eventPaymentCountByEventID = `-- name: EventPaymentCountByEventID :one
SELECT COUNT(*) FROM (SELECT DISTINCT ep.id FROM event_payment_receipts as ep
INNER JOIN event_registrations as er ON er.event_payment_receipt_id = ep.id
INNER JOIN clubs as c ON c.id = er.club_id
WHERE er.event_id = $1) AS count
`

const eventPaymentCountByStatus = `-- name: EventPaymentCountByStatus :one
SELECT COUNT(*) FROM (SELECT DISTINCT ep.id FROM event_payment_receipts as ep
INNER JOIN event_registrations as er ON er.event_payment_receipt_id = ep.id
INNER JOIN clubs as c ON c.id = er.club_id
WHERE ep.status = $1) AS count
`

const eventPaymentCountByClubID = `-- name: EventPaymentCountByClubID :one
SELECT COUNT(*) FROM (SELECT DISTINCT ep.id FROM event_payment_receipts as ep
INNER JOIN event_registrations as er ON er.event_payment_receipt_id = ep.id
INNER JOIN clubs as c ON c.id = er.club_id
WHERE er.club_id = any($1)) AS count
`

const eventPaymentCountByDate = `-- name: EventPaymentCountByDate :one
SELECT COUNT(*) FROM (SELECT DISTINCT ep.id FROM event_payment_receipts as ep
INNER JOIN event_registrations as er ON er.event_payment_receipt_id = ep.id
INNER JOIN clubs as c ON c.id = er.club_id
WHERE ep.created_at BETWEEN DATE($1) AND DATE($2)) AS count
`

const eventPaymentCountByEventIDAndStatus = `-- name: EventPaymentCountByEventIDAndStatus :one
SELECT COUNT(*) FROM (SELECT DISTINCT ep.id FROM event_payment_receipts as ep
INNER JOIN event_registrations as er ON er.event_payment_receipt_id = ep.id
INNER JOIN clubs as c ON c.id = er.club_id
WHERE er.event_id = $1 AND ep.status = $2) AS count
`

const eventPaymentCountByEventIDAndClubID = `-- name: EventPaymentCountByEventIDAndClubID :one
SELECT COUNT(*) FROM (SELECT DISTINCT ep.id FROM event_payment_receipts as ep
INNER JOIN event_registrations as er ON er.event_payment_receipt_id = ep.id
INNER JOIN clubs as c ON c.id = er.club_id
WHERE er.event_id = $1 AND er.club_id = any($2)) AS count
`

const eventPaymentCountByEventIDAndDate = `-- name: EventPaymentCountByEventIDAndDate :one
SELECT COUNT(*) FROM (SELECT DISTINCT ep.id FROM event_payment_receipts as ep
INNER JOIN event_registrations as er ON er.event_payment_receipt_id = ep.id
INNER JOIN clubs as c ON c.id = er.club_id
WHERE er.event_id = $1 AND ep.created_at BETWEEN DATE($2) AND DATE($3)) AS count
`

const eventPaymentCountByStatusAndClubID = `-- name: EventPaymentCountByStatusAndClubID :one
SELECT COUNT(*) FROM (SELECT DISTINCT ep.id FROM event_payment_receipts as ep
INNER JOIN event_registrations as er ON er.event_payment_receipt_id = ep.id
INNER JOIN clubs as c ON c.id = er.club_id
WHERE ep.status = $1 AND er.club_id = any($2)) AS count
`

const eventPaymentCountByStatusAndDate = `-- name: EventPaymentCountByStatusAndClubID :one
SELECT COUNT(*) FROM (SELECT DISTINCT ep.id FROM event_payment_receipts as ep
INNER JOIN event_registrations as er ON er.event_payment_receipt_id = ep.id
INNER JOIN clubs as c ON c.id = er.club_id
WHERE ep.status = $1 AND ep.created_at BETWEEN DATE($2) AND DATE($3)) AS count
`

const eventPaymentCountByClubIDAndDate = `-- name: EventPaymentCountByClubIDAndClubID :one
SELECT COUNT(*) FROM (SELECT DISTINCT ep.id FROM event_payment_receipts as ep
INNER JOIN event_registrations as er ON er.event_payment_receipt_id = ep.id
INNER JOIN clubs as c ON c.id = er.club_id
WHERE er.club_id = $1 AND ep.created_at BETWEEN DATE($2) AND DATE($3)) AS count
`

const eventPaymentCountByEventIDStatusAndClubID = `-- name: EventPaymentCountByEventIDStatusAndClubID :one
SELECT COUNT(*) FROM (SELECT DISTINCT ep.id FROM event_payment_receipts as ep
INNER JOIN event_registrations as er ON er.event_payment_receipt_id = ep.id
INNER JOIN clubs as c ON c.id = er.club_id
WHERE er.event_id = $1 AND ep.status = $2 AND er.club_id = any($3)) AS count
`

const eventPaymentCountByEventIDStatusAndDate = `-- name: EventPaymentCountByEventIDStatusAndDate :one
SELECT COUNT(*) FROM (SELECT DISTINCT ep.id FROM event_payment_receipts as ep
INNER JOIN event_registrations as er ON er.event_payment_receipt_id = ep.id
INNER JOIN clubs as c ON c.id = er.club_id
WHERE er.event_id = $1 AND ep.status = $2 AND ep.created_at BETWEEN DATE($3) AND DATE($4)) AS count
`

const eventPaymentCountByEventIDClubIDAndDate = `-- name: EventPaymentCountByEventIDStatusAndDate :one
SELECT COUNT(*) FROM (SELECT DISTINCT ep.id FROM event_payment_receipts as ep
INNER JOIN event_registrations as er ON er.event_payment_receipt_id = ep.id
INNER JOIN clubs as c ON c.id = er.club_id
WHERE er.event_id = $1 AND er.club_id = any($2) AND ep.created_at BETWEEN DATE($3) AND DATE($4)) AS count
`

const eventPaymentCountByStatusClubIDAndDate = `-- name: EventPaymentCountByStatusClubIDAndDate :one
SELECT COUNT(*) FROM (SELECT DISTINCT ep.id FROM event_payment_receipts as ep
INNER JOIN event_registrations as er ON er.event_payment_receipt_id = ep.id
INNER JOIN clubs as c ON c.id = er.club_id
WHERE ep.status = $1 AND er.club_id = any($2) AND BETWEEN DATE($3) AND DATE($4)) AS count
`

const eventPaymentCountByEventIDFilterAll = `-- name: EventPaymentCountByEventIDFilterAll :one
SELECT COUNT(*) FROM (SELECT DISTINCT ep.id FROM event_payment_receipts as ep
INNER JOIN event_registrations as er ON er.event_payment_receipt_id = ep.id
INNER JOIN clubs as c ON c.id = er.club_id
WHERE er.event_id = $1 AND ep.status AND er.club_id = any($2) AND ep.created_at BETWEEN DATE($3) AND DATE($4)) AS count
`

func (r *RawRepository) EventPaymentCountByEventID(ctx context.Context, arg queryPaymentParams) (int64, error) {
	var row *sql.Row
	status := r.setQueryStatus(arg)
	switch status {
	default:
		row = r.db.QueryRowContext(ctx, eventPaymentCountByAll)
	case "event":
		row = r.db.QueryRowContext(ctx, eventPaymentCountByEventID, arg.EventID)
	case "status":
		row = r.db.QueryRowContext(ctx, eventPaymentCountByStatus, arg.Status)
	case "club_id":
		row = r.db.QueryRowContext(ctx, eventPaymentCountByClubID, pq.Array(arg.ClubID))
	case "date":
		row = r.db.QueryRowContext(ctx, eventPaymentCountByDate, arg.Start, arg.End)
	case "event-status":
		row = r.db.QueryRowContext(ctx, eventPaymentCountByEventIDAndStatus, arg.EventID, arg.Status)
	case "event-club":
		row = r.db.QueryRowContext(ctx, eventPaymentCountByEventIDAndClubID, arg.EventID, pq.Array(arg.ClubID))
	case "event-date":
		row = r.db.QueryRowContext(ctx, eventPaymentCountByEventIDAndDate, arg.EventID, arg.Start, arg.End)
	case "status-club_id":
		row = r.db.QueryRowContext(ctx, eventPaymentCountByStatusAndClubID, arg.Status, pq.Array(arg.ClubID))
	case "status-date":
		row = r.db.QueryRowContext(ctx, eventPaymentCountByStatusAndDate, arg.Status, arg.Start, arg.End)
	case "club_id-date":
		row = r.db.QueryRowContext(ctx, eventPaymentCountByClubIDAndDate, pq.Array(arg.ClubID), arg.Start, arg.End)
	case "event-status-club":
		row = r.db.QueryRowContext(ctx, eventPaymentCountByEventIDStatusAndClubID, arg.EventID, arg.Status, pq.Array(arg.ClubID))
	case "event-status-date":
		row = r.db.QueryRowContext(ctx, eventPaymentCountByEventIDStatusAndDate, arg.EventID, arg.Status, arg.Start, arg.End)
	case "event-club-date":
		row = r.db.QueryRowContext(ctx, eventPaymentCountByEventIDClubIDAndDate, arg.EventID, pq.Array(arg.ClubID), arg.Start, arg.End)
	case "status-club_id-date":
		row = r.db.QueryRowContext(ctx, eventPaymentCountByStatusClubIDAndDate, arg.Status, pq.Array(arg.ClubID), arg.Start, arg.End)
	case "event-status-club_id-date":
		row = r.db.QueryRowContext(ctx, eventPaymentCountByEventIDFilterAll, arg.EventID, arg.Status, pq.Array(arg.ClubID), arg.Start, arg.End)
	}
	var count int64
	err := row.Scan(&count)
	return count, err
}

const eventRegistrationFetchCartDetails = `-- name: EventRegistrationFetchCartDetails :many
SELECT er.id, c.name as class_name, ce.price, array(
    SELECT u.name FROM event_participants as ep
    INNER JOIN users as u ON u.id = ep.user_id
    WHERE ep.event_registration_id = er.id
) as participants FROM event_registrations as er
INNER JOIN class_events as ce ON ce.id = er.class_event_id
INNER JOIN classes as c ON c.id = ce.class_id
INNER JOIN clubs as cl ON cl.id = er.club_id
WHERE er.status = 'pending' AND cl.user_id = $1 AND er.event_id = $2
`

const eventRegistrationFetchCartDetailsWithoutUserID = `-- name: EventRegistrationFetchCartDetailsWithoutUserID :many
SELECT er.id, c.name as class_name, ce.price, array(
    SELECT u.name FROM event_participants as ep
    INNER JOIN users as u ON u.id = ep.user_id
    WHERE ep.event_registration_id = er.id
) as participants FROM event_registrations as er
INNER JOIN class_events as ce ON ce.id = er.class_event_id
INNER JOIN classes as c ON c.id = ce.class_id
INNER JOIN clubs as cl ON cl.id = er.club_id
WHERE er.event_id = $1
`

type EventRegistrationFetchCartDetailsParams struct {
	UserID  uuid.UUID `json:"club_id"`
	EventID uuid.UUID `json:"event_id"`
}

type EventRegistrationFetchCartDetailsRow struct {
	ID           uuid.UUID `json:"id"`
	ClassName    string    `json:"class_name"`
	Price        int32     `json:"price"`
	Participants []string  `json:"participants"`
}

func (r *RawRepository) EventRegistrationFetchCartDetails(ctx context.Context, arg EventRegistrationFetchCartDetailsParams) ([]EventRegistrationFetchCartDetailsRow, error) {
	var rows *sql.Rows
	if arg.UserID.ID() != 0 {
		result, err := r.db.QueryContext(ctx, eventRegistrationFetchCartDetails, arg.UserID, arg.EventID)
		if err != nil {
			return nil, err
		}
		rows = result
	} else {
		result, err := r.db.QueryContext(ctx, eventRegistrationFetchCartDetailsWithoutUserID, arg.EventID)
		if err != nil {
			return nil, err
		}
		rows = result
	}

	defer rows.Close()
	items := []EventRegistrationFetchCartDetailsRow{}
	for rows.Next() {
		var i EventRegistrationFetchCartDetailsRow
		if err := rows.Scan(
			&i.ID,
			&i.ClassName,
			&i.Price,
			pq.Array(&i.Participants),
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

const eventPaymentSumAll = `-- name: EventPaymentSumAll :one
SELECT COALESCE(COALESCE(SUM(ce.price),0),0) FROM event_registrations AS er
INNER JOIN class_events AS ce ON ce.id = er.class_event_id
INNER JOIN event_payment_receipts as ep ON ep.id = er.event_payment_receipt_id
INNER JOIN classes as c ON c.id = ce.class_id
`

const eventPaymentSumByEventID = `-- name: EventPaymentSumByEventID :one
SELECT COALESCE(SUM(ce.price),0) FROM event_registrations AS er
INNER JOIN class_events AS ce ON ce.id = er.class_event_id
INNER JOIN event_payment_receipts as ep ON ep.id = er.event_payment_receipt_id
INNER JOIN classes as c ON c.id = ce.class_id
WHERE er.event_id = $1
`

const eventPaymentSumByStatus = `-- name: EventPaymentSumByStatus :one
SELECT COALESCE(SUM(ce.price),0) FROM event_registrations AS er
INNER JOIN class_events AS ce ON ce.id = er.class_event_id
INNER JOIN event_payment_receipts as ep ON ep.id = er.event_payment_receipt_id
INNER JOIN classes as c ON c.id = ce.class_id
WHERE ep.status = $1
`

const eventPaymentSumByFilterAll = `-- name: EventPaymentSumByFilterAll :one
SELECT COALESCE(SUM(ce.price),0) FROM event_registrations AS er
INNER JOIN class_events AS ce ON ce.id = er.class_event_id
INNER JOIN event_payment_receipts as ep ON ep.id = er.event_payment_receipt_id
INNER JOIN classes as c ON c.id = ce.class_id
WHERE er.event_id = $1 AND ep.status = $2
`

type sumAllParams struct {
	EventID string                `json:"event_id"`
	Status  db.EventReceiptStatus `json:"status"`
}

func (r *RawRepository) setQuerySumStatus(arg sumAllParams) string {
	var result []string
	if arg.EventID != "" {
		result = append(result, "event")
	}
	if arg.Status != "" {
		result = append(result, "status")
	}
	return strings.Join(result, "-")
}

func (r *RawRepository) EventPaymentSumAll(ctx context.Context, arg sumAllParams) (int64, error) {
	var row *sql.Row
	status := r.setQuerySumStatus(arg)
	switch status {
	default:
		row = r.db.QueryRowContext(ctx, eventPaymentSumAll)
	case "event":
		row = r.db.QueryRowContext(ctx, eventPaymentSumByEventID, arg.EventID)
	case "status":
		row = r.db.QueryRowContext(ctx, eventPaymentSumByStatus, arg.Status)
	case "event-status":
		row = r.db.QueryRowContext(ctx, eventPaymentSumByFilterAll, arg.EventID, arg.Status)
	}
	var sum int64
	err := row.Scan(&sum)
	return sum, err
}
