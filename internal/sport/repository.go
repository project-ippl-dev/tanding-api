package sport

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/project-ippl-dev/tanding-api/internal/db"
)

type RawRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) RawRepository {
	return RawRepository{db: db}
}

const sportFetchAll = `-- name: SportFetchAll :many
SELECT id, name, description, type, thumbnail FROM sports WHERE deleted_at IS NULL AND name LIKE $1 ORDER BY name LIMIT $2 OFFSET $3
`

const sportFetchByCategory = `-- name: SportFetchByCategory :many
SELECT id, name, description, type, thumbnail FROM sports WHERE deleted_at IS NULL AND type = $1 AND name LIKE $2 
ORDER BY name LIMIT $3 OFFSET $4
`

type fetchAllParams struct {
	Name     string `json:"name"`
	Limit    int32  `json:"limit"`
	Offset   int32  `json:"offset"`
	Category string `json:"category"`
}

type FetchAllRow struct {
	ID          uuid.UUID    `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Type        db.SportType `json:"type"`
	Thumbnail   string       `json:"thumbnail"`
}

func (r *RawRepository) SportFetchAll(ctx context.Context, arg fetchAllParams) ([]FetchAllRow, error) {
	var rows *sql.Rows
	if arg.Category == "" {
		result, err := r.db.QueryContext(ctx, sportFetchAll, arg.Name, arg.Limit, arg.Offset)
		if err != nil {
			return nil, err
		}
		rows = result
	} else {
		result, err := r.db.QueryContext(ctx, sportFetchByCategory, arg.Category, arg.Name, arg.Limit, arg.Offset)
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
			&i.Description,
			&i.Type,
			&i.Thumbnail,
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

const sportCountAll = `-- name: SportCountAll :one
SELECT COUNT(id) FROM sports WHERE deleted_at IS NULL AND name LIKE $1
`

const sportCountByCategory = `-- name: SportCountByCategory :one
SELECT COUNT(id) FROM sports WHERE deleted_at IS NULL AND type = $1 AND name LIKE $2;
`

type countAllParams struct {
	Name     string `json:"name"`
	Category string `json:"category"`
}

func (r *RawRepository) SportCountAll(ctx context.Context, args countAllParams) (int64, error) {
	var row *sql.Row
	if args.Category == "" {
		row = r.db.QueryRowContext(ctx, sportCountAll, args.Name)
	} else {
		row = r.db.QueryRowContext(ctx, sportCountByCategory, args.Category, args.Name)
	}
	var count int64
	err := row.Scan(&count)
	return count, err
}
