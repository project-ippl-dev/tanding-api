package accomplishment

//go:generate mockgen -source=./usecase.go -destination=../../mocks/accomplishment/usecase_mock.go

import (
	"context"
	"fmt"
	"time"

	"database/sql"
	"github.com/google/uuid"

	"github.com/project-ippl-dev/tanding-api/internal/db"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
)

type Usecase interface {
	Store(ctx context.Context, req Request, userID string) error
	FetchAll(ctx context.Context, page int32, pageSize int32, userID string) (tools.Pagination, error)
	Update(ctx context.Context, req Request, userID string, accomplishmentID int64) error
	Delete(ctx context.Context, userID string, accomplishmentID int64) error
}
type usecase struct {
	repository *db.Queries
}

func NewUsecase(repository *db.Queries) Usecase {
	return &usecase{repository: repository}
}

func (u usecase) Store(ctx context.Context, req Request, userID string) error {
	return u.repository.AccomplishmentCreate(ctx, db.AccomplishmentCreateParams{
		UserID:   uuid.MustParse(userID),
		Title:    req.Title,
		Level:    req.Level,
		Ranking:  req.Ranking,
		Category: req.Category,
		Sport:    req.Sport,
		Description: sql.NullString{
			String: req.Description,
			Valid:  true,
		},
		FileUrl: req.FileURL,
		Month:   req.Month,
		Year:    req.Year,
	})
}

func (u usecase) FetchAll(ctx context.Context, page int32, pageSize int32, userID string) (tools.Pagination, error) {
	skip := tools.PaginationSkip(page, pageSize)
	accomplishments, err := u.repository.AccomplishmentFetchAll(ctx, db.AccomplishmentFetchAllParams{
		UserID: uuid.MustParse(userID),
		Limit:  pageSize,
		Offset: skip,
	})
	if err != nil {
		return tools.Pagination{}, fmt.Errorf("error in fetch accomplishment : %s", err.Error())
	}

	result := []Response{}

	var date string
	for _, accomplishment := range accomplishments {
		if accomplishment.Month != 0 {
			month := time.Month(accomplishment.Month)
			date = fmt.Sprintf("%s %d", month.String(), accomplishment.Year)
		} else {
			date = fmt.Sprintf("%d", accomplishment.Year)
		}
		data := Response{
			Title:       accomplishment.Title,
			Level:       accomplishment.Level,
			Ranking:     accomplishment.Ranking,
			Category:    accomplishment.Category,
			Sport:       accomplishment.Sport,
			Description: accomplishment.Description.String,
			FileURL:     accomplishment.FileUrl,
			Date:        date,
		}
		result = append(result, data)
	}

	count, err := u.repository.AccomplishmentCountAll(ctx, uuid.MustParse(userID))
	if err != nil {
		return tools.Pagination{}, fmt.Errorf("error in count accomplishment : %s", err.Error())
	}
	return tools.Pagination{
		TotalItem: count,
		PageSize:  pageSize,
		Page:      page,
		Data:      result,
	}, nil
}

func (u usecase) Update(ctx context.Context, req Request, userID string, accomplishmentID int64) error {
	accomplishID, err := u.repository.AccomplishmentCheckOne(ctx, db.AccomplishmentCheckOneParams{
		ID:     accomplishmentID,
		UserID: uuid.MustParse(userID),
	})
	if err != nil {
		return fmt.Errorf("error in check one accomplishment : %s", err.Error())
	}
	return u.repository.AccomplishmentUpdate(ctx, db.AccomplishmentUpdateParams{
		Title:    req.Title,
		Level:    req.Level,
		Ranking:  req.Ranking,
		Category: req.Category,
		Sport:    req.Sport,
		Description: sql.NullString{
			String: req.Description,
			Valid:  true,
		},
		FileUrl: req.FileURL,
		Month:   req.Month,
		Year:    req.Year,
		ID:      accomplishID,
	})
}

func (u usecase) Delete(ctx context.Context, userID string, accomplishmentID int64) error {
	accomplishID, err := u.repository.AccomplishmentCheckOne(ctx, db.AccomplishmentCheckOneParams{
		ID:     accomplishmentID,
		UserID: uuid.MustParse(userID),
	})
	if err != nil {
		return fmt.Errorf("error in check one accomplishment : %s", err.Error())
	}
	return u.repository.AccomplishmentDelete(ctx, accomplishID)
}
