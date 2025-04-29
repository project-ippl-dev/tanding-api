package event

import (
	"context"
	"database/sql"
	"github.com/dytlan/tanding-api/internal/db"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"strings"
	"time"
)

type RawRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) RawRepository {
	return RawRepository{db: db}
}

const fetchInfinite = `-- name: EventFetchInfinite :many
SELECT e.id, e.user_id, u.name as user_name, e.type, e.name, e.description, e.prize_pool, e.location, e.thumbnail, e.start_date, e.end_date,
       e.deadline, e.sport_id, s.name as sport_name, e.order_number, e.quota, e.remark, e.open, e.province, e.city, u.photo as user_image
FROM events as e
         INNER JOIN users as u ON e.user_id = u.id
         INNER JOIN sports as s ON e.sport_id = s.id
WHERE e.deleted_at IS NULL AND e.remark != 'unconfirmed' AND e.remark != 'rejected' AND e.status = TRUE AND e.order_number < $1 AND e.name LIKE $2
ORDER BY e.order_number DESC,
CASE e.remark
	WHEN 'ongoing' THEN 1
	WHEN 'closed' THEN 2
	WHEN 'open' THEN 3
	WHEN 'soon' THEN 4
	WHEN 'done' THEN 5
	ELSE 6
END
LIMIT $3
`

const fetchInfiniteBySportID = `-- name: EventFetchInfiniteBySportID :many
SELECT e.id, e.user_id, u.name as user_name, e.type, e.name, e.description, e.prize_pool, e.location, e.thumbnail, e.start_date, e.end_date,
       e.deadline, e.sport_id, s.name as sport_name, e.order_number, e.quota, e.remark, e.open, e.province, e.city, u.photo as user_image
FROM events as e
         INNER JOIN users as u ON e.user_id = u.id
         INNER JOIN sports as s ON e.sport_id = s.id
WHERE e.deleted_at IS NULL AND e.remark != 'unconfirmed' AND e.remark != 'rejected' AND e.status = TRUE AND e.order_number < $1 AND e.name LIKE $2 AND e.sport_id = $3
ORDER BY e.order_number DESC,
CASE e.remark
	WHEN 'ongoing' THEN 1
	WHEN 'closed' THEN 2
	WHEN 'open' THEN 3
	WHEN 'soon' THEN 4
	WHEN 'done' THEN 5
	ELSE 6
END
LIMIT $4
`

const fetchInfiniteByCategory = `-- name: EventFetchInfiniteByCategory :many
SELECT e.id, e.user_id, u.name as user_name, e.type, e.name, e.description, e.prize_pool, e.location, e.thumbnail, e.start_date, e.end_date,
       e.deadline, e.sport_id, s.name as sport_name, e.order_number, e.quota, e.remark, e.open, e.province, e.city, u.photo as user_image
FROM events as e
         INNER JOIN users as u ON e.user_id = u.id
         INNER JOIN sports as s ON e.sport_id = s.id
WHERE e.deleted_at IS NULL AND e.remark != 'unconfirmed' AND e.remark != 'rejected' AND e.status = TRUE AND e.order_number < $1 AND e.name LIKE $2 AND s.type = $3
ORDER BY e.order_number DESC,
CASE e.remark
	WHEN 'ongoing' THEN 1
	WHEN 'closed' THEN 2
	WHEN 'open' THEN 3
	WHEN 'soon' THEN 4
	WHEN 'done' THEN 5
	ELSE 6
END
LIMIT $4
`

const fetchInfiniteByRemark = `-- name: EventFetchInfiniteByRemark :many
SELECT e.id, e.user_id, u.name as user_name, e.type, e.name, e.description, e.prize_pool, e.location, e.thumbnail, e.start_date, e.end_date,
       e.deadline, e.sport_id, s.name as sport_name, e.order_number, e.quota, e.remark, e.open, e.province, e.city, u.photo as user_image
FROM events as e
         INNER JOIN users as u ON e.user_id = u.id
         INNER JOIN sports as s ON e.sport_id = s.id
WHERE e.deleted_at IS NULL AND e.remark != 'unconfirmed' AND e.remark != 'rejected' AND e.status = TRUE AND e.order_number < $1 AND e.name LIKE $2 AND e.remark = $3
ORDER BY e.order_number DESC,
CASE e.remark
	WHEN 'ongoing' THEN 1
	WHEN 'closed' THEN 2
	WHEN 'open' THEN 3
	WHEN 'soon' THEN 4
	WHEN 'done' THEN 5
	ELSE 6
END
LIMIT $4
`

const fetchInfiniteBySportIDAndCategory = `-- name: EventFetchInfiniteByCategory :many
SELECT e.id, e.user_id, u.name as user_name, e.type, e.name, e.description, e.prize_pool, e.location, e.thumbnail, e.start_date, e.end_date,
       e.deadline, e.sport_id, s.name as sport_name, e.order_number, e.quota, e.remark, e.open, e.province, e.city, u.photo as user_image
FROM events as e
         INNER JOIN users as u ON e.user_id = u.id
         INNER JOIN sports as s ON e.sport_id = s.id
WHERE e.deleted_at IS NULL AND e.remark != 'unconfirmed' AND e.remark != 'rejected' AND e.status = TRUE AND e.order_number < $1 AND e.name LIKE $2 AND e.sport_id = $3 AND s.type = $4
ORDER BY e.order_number DESC,
CASE e.remark
	WHEN 'ongoing' THEN 1
	WHEN 'closed' THEN 2
	WHEN 'open' THEN 3
	WHEN 'soon' THEN 4
	WHEN 'done' THEN 5
	ELSE 6
END
LIMIT $5
`

const fetchInfiniteBySportIDAndRemark = `-- name: EventFetchInfiniteBySportIDAndRemark :many
SELECT e.id, e.user_id, u.name as user_name, e.type, e.name, e.description, e.prize_pool, e.location, e.thumbnail, e.start_date, e.end_date,
       e.deadline, e.sport_id, s.name as sport_name, e.order_number, e.quota, e.remark, e.open, e.province, e.city, u.photo as user_image
FROM events as e
         INNER JOIN users as u ON e.user_id = u.id
         INNER JOIN sports as s ON e.sport_id = s.id
WHERE e.deleted_at IS NULL AND e.remark != 'unconfirmed' AND e.remark != 'rejected' AND e.status = TRUE AND e.order_number < $1 AND e.name LIKE $2 AND e.sport_id = $3 AND e.remark = $4
ORDER BY e.order_number DESC,
CASE e.remark
	WHEN 'ongoing' THEN 1
	WHEN 'closed' THEN 2
	WHEN 'open' THEN 3
	WHEN 'soon' THEN 4
	WHEN 'done' THEN 5
	ELSE 6
END
LIMIT $5
`

const fetchInfiniteByCategoryAndRemark = `-- name: EventFetchInfiniteByCategoryAndRemark :many
SELECT e.id, e.user_id, u.name as user_name, e.type, e.name, e.description, e.prize_pool, e.location, e.thumbnail, e.start_date, e.end_date,
       e.deadline, e.sport_id, s.name as sport_name, e.order_number, e.quota, e.remark, e.open, e.province, e.city, u.photo as user_image
FROM events as e
         INNER JOIN users as u ON e.user_id = u.id
         INNER JOIN sports as s ON e.sport_id = s.id
WHERE e.deleted_at IS NULL AND e.remark != 'unconfirmed' AND e.remark != 'rejected' AND e.status = TRUE AND e.order_number < $1 AND e.name LIKE $2 AND e.sport_id = $3 AND e.remark = $4
ORDER BY e.order_number DESC,
CASE e.remark
	WHEN 'ongoing' THEN 1
	WHEN 'closed' THEN 2
	WHEN 'open' THEN 3
	WHEN 'soon' THEN 4
	WHEN 'done' THEN 5
	ELSE 6
END
LIMIT $5
`

const fetchInfiniteByFilterall = `-- name: EventFetchInfiniteByFilterAll :many
SELECT e.id, e.user_id, u.name as user_name, e.type, e.name, e.description, e.prize_pool, e.location, e.thumbnail, e.start_date, e.end_date,
       e.deadline, e.sport_id, s.name as sport_name, e.order_number, e.quota, e.remark, e.open, e.province, e.city, u.photo as user_image
FROM events as e
         INNER JOIN users as u ON e.user_id = u.id
         INNER JOIN sports as s ON e.sport_id = s.id
WHERE e.deleted_at IS NULL AND e.remark != 'unconfirmed' AND e.remark != 'rejected' AND e.status = TRUE AND e.order_number < $1 
  AND e.name LIKE $2 AND e.sport_id = $3 AND s.type = $4 AND e.remark = $5
ORDER BY e.order_number DESC,
CASE e.remark
	WHEN 'ongoing' THEN 1
	WHEN 'closed' THEN 2
	WHEN 'open' THEN 3
	WHEN 'soon' THEN 4
	WHEN 'done' THEN 5
	ELSE 6
END
LIMIT $6
`

type fetchQueryParams struct {
	SportID  string        `json:"sport_id"`
	Name     string        `json:"name"`
	Category db.SportType  `json:"category"`
	Remark   db.RemarkType `json:"remark"`
}

type fetchInfiniteParams struct {
	fetchQueryParams
	OrderNumber int64 `json:"order_number"`
	Limit       int32 `json:"limit"`
}

type fetchInfiniteRow struct {
	ID          uuid.UUID     `json:"id"`
	UserID      uuid.UUID     `json:"user_id"`
	UserName    string        `json:"user_name"`
	Type        db.EventType  `json:"type"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	PrizePool   string        `json:"prize_pool"`
	Location    string        `json:"location"`
	Province    string        `json:"province"`
	City        string        `json:"city"`
	Thumbnail   string        `json:"thumbnail"`
	StartDate   time.Time     `json:"start_date"`
	EndDate     time.Time     `json:"end_date"`
	Deadline    time.Time     `json:"deadline"`
	SportID     uuid.UUID     `json:"sport_id"`
	SportName   string        `json:"sport_name"`
	OrderNumber int64         `json:"order_number"`
	Quota       int32         `json:"quota"`
	Remark      db.RemarkType `json:"remark"`
	Open        time.Time     `json:"open"`
	UserImage   string        `json:"user_image"`
}

func (r *RawRepository) setQueryStatus(args fetchQueryParams) string {
	var results []string

	if args.SportID != "" {
		results = append(results, "sport")
	}
	if args.Category != "" {
		results = append(results, "category")
	}
	if args.Remark != "" {
		results = append(results, "remark")
	}

	return strings.Join(results, "-")
}

func (r *RawRepository) EventFetchInfinite(ctx context.Context, args fetchInfiniteParams) ([]fetchInfiniteRow, error) {
	var rows *sql.Rows
	queryStatus := r.setQueryStatus(args.fetchQueryParams)
	switch queryStatus {
	default:
		result, err := r.db.QueryContext(ctx, fetchInfinite, args.OrderNumber, args.Name, args.Limit)
		if err != nil {
			return nil, err
		}
		rows = result
		break
	case "sport":
		result, err := r.db.QueryContext(ctx, fetchInfiniteBySportID, args.OrderNumber, args.Name, args.SportID, args.Limit)
		if err != nil {
			return nil, err
		}
		rows = result
		break
	case "category":
		result, err := r.db.QueryContext(ctx, fetchInfiniteByCategory, args.OrderNumber, args.Name, args.Category, args.Limit)
		if err != nil {
			return nil, err
		}
		rows = result
		break
	case "remark":
		result, err := r.db.QueryContext(ctx, fetchInfiniteByRemark, args.OrderNumber, args.Name, args.Remark, args.Limit)
		if err != nil {
			return nil, err
		}
		rows = result
		break
	case "sport-category":
		result, err := r.db.QueryContext(ctx, fetchInfiniteBySportIDAndCategory, args.OrderNumber, args.Name, args.SportID, args.Category, args.Limit)
		if err != nil {
			return nil, err
		}
		rows = result
		break
	case "sport-remark":
		result, err := r.db.QueryContext(ctx, fetchInfiniteBySportIDAndRemark, args.OrderNumber, args.Name, args.SportID, args.Remark, args.Limit)
		if err != nil {
			return nil, err
		}
		rows = result
		break
	case "category-remark":
		result, err := r.db.QueryContext(ctx, fetchInfiniteByCategoryAndRemark, args.OrderNumber, args.Name, args.Category, args.Remark, args.Limit)
		if err != nil {
			return nil, err
		}
		rows = result
		break
	case "sport-category-remark":
		result, err := r.db.QueryContext(ctx, fetchInfiniteByFilterall, args.OrderNumber, args.Name, args.SportID, args.Category, args.Remark, args.Limit)
		if err != nil {
			return nil, err
		}
		rows = result
		break
	}
	defer rows.Close()
	items := []fetchInfiniteRow{}
	for rows.Next() {
		var i fetchInfiniteRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.UserName,
			&i.Type,
			&i.Name,
			&i.Description,
			&i.PrizePool,
			&i.Location,
			&i.Thumbnail,
			&i.StartDate,
			&i.EndDate,
			&i.Deadline,
			&i.SportID,
			&i.SportName,
			&i.OrderNumber,
			&i.Quota,
			&i.Remark,
			&i.Open,
			&i.Province,
			&i.City,
			&i.UserImage,
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

const fetchLatestOrder = `-- name: EventFetchLatestOrder :one
SELECT order_number FROM events WHERE remark != 'unconfirmed' AND name LIKE $1 ORDER BY order_number DESC LIMIT 1
`

const fetchLatestOrderBySportID = `-- name: EventFetchLatestOrderBySportID :one
SELECT e.order_number FROM events as e
INNER JOIN sports as s ON e.sport_id = s.id
WHERE e.remark != 'unconfirmed' AND e.name LIKE $1 AND e.sport_id = $2
ORDER BY order_number DESC LIMIT 1
`

const fetchLatestOrderByCategory = `-- name: EventFetchLatestOrderByCategory :one
SELECT e.order_number FROM events as e
INNER JOIN sports as s ON e.sport_id = s.id
WHERE e.remark != 'unconfirmed' AND e.name LIKE $1 AND s.type = $2
ORDER BY order_number DESC LIMIT 1
`

const fetchLatestOrderByRemark = `-- name: EventFetchLatestOrderByRemark :one
SELECT e.order_number FROM events as e
INNER JOIN sports as s ON e.sport_id = s.id
WHERE e.remark != 'unconfirmed' AND e.name LIKE $1 AND e.remark = $2
ORDER BY order_number DESC LIMIT 1
`

const fetchLatestOrderBySportIDAndCategory = `-- name: EventFetchLatestOrderBySportIDAndCategory :one
SELECT e.order_number FROM events as e
INNER JOIN sports as s ON e.sport_id = s.id
WHERE e.remark != 'unconfirmed' AND e.name LIKE $1 AND e.sport_id = $2 AND s.type $3 
ORDER BY order_number DESC LIMIT 1
`

const fetchLatestOrderBySportIDAndRemark = `-- name: EventFetchLatestOrderBySportIDAndRemark :one
SELECT e.order_number FROM events as e
INNER JOIN sports as s ON e.sport_id = s.id
WHERE e.remark != 'unconfirmed' AND e.name LIKE $1 AND e.sport_id = $2 AND e.remark $3 
ORDER BY order_number DESC LIMIT 1
`

const fetchLatestOrderByCategoryAndRemark = `-- name: EventFetchLatestOrderBySportIDAndRemark :one
SELECT e.order_number FROM events as e
INNER JOIN sports as s ON e.sport_id = s.id
WHERE e.remark != 'unconfirmed' AND e.name LIKE $1 AND s.type = $2 AND e.remark $3 
ORDER BY order_number DESC LIMIT 1
`

const fetchLatestOrderByFilterAll = `-- name: EventFetchLatestOrderByFilterAll :one
SELECT e.order_number FROM events as e
INNER JOIN sports as s ON e.sport_id = s.id
WHERE e.remark != 'unconfirmed' AND e.name LIKE $1 AND e.sport_id = $2 AND s.type = $3 AND e.remark $4 
ORDER BY order_number DESC LIMIT 1
`

func (r *RawRepository) EventFetchLatestOrder(ctx context.Context, args fetchQueryParams) (int64, error) {
	var row *sql.Row
	queryStatus := r.setQueryStatus(args)
	switch queryStatus {
	default:
		row = r.db.QueryRowContext(ctx, fetchLatestOrder, args.Name)
		break
	case "sport":
		row = r.db.QueryRowContext(ctx, fetchLatestOrderBySportID, args.Name, args.SportID)
		break
	case "category":
		row = r.db.QueryRowContext(ctx, fetchLatestOrderByCategory, args.Name, args.Category)
		break
	case "remark":
		row = r.db.QueryRowContext(ctx, fetchLatestOrderByRemark, args.Name, args.Remark)
		break
	case "sport-category":
		row = r.db.QueryRowContext(ctx, fetchLatestOrderBySportIDAndCategory, args.Name, args.SportID, args.Category)
		break
	case "sport-remark":
		row = r.db.QueryRowContext(ctx, fetchLatestOrderBySportIDAndRemark, args.Name, args.SportID, args.Remark)
		break
	case "category-remark":
		row = r.db.QueryRowContext(ctx, fetchLatestOrderByCategoryAndRemark, args.Name, args.Category, args.Remark)
		break
	case "sport-category-remark":
		row = r.db.QueryRowContext(ctx, fetchLatestOrderByFilterAll, args.Name, args.SportID, args.Category, args.Remark)
		break
	}
	var order_number int64
	err := row.Scan(&order_number)
	return order_number, err
}

const rankFetchByClassEventID = `-- name: RankFetchByClassEventID :many
SELECT r.id, r.rank, r.point, c.name as club_name ,array(
    SELECT u.name FROM event_participants as ep
    INNER JOIN users as u ON u.id = ep.user_id
    WHERE ep.event_registration_id = er.id
) as participants FROM ranks as r
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
