package sport

//go:generate mockgen -source=./usecase.go -destination=../../mocks/sport/usecase_mock.go

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/project-ippl-dev/tanding-api/internal/db"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
)

type Usecase interface {
	Store(ctx context.Context, req Request) error
	FetchAll(ctx context.Context, page int32, pageSize int32, queries FetchAllQueryParams) (tools.Pagination, error)
	Update(ctx context.Context, req Request, sportID string) error
	Delete(ctx context.Context, sportID string) error
}

type usecase struct {
	repository    *db.Queries
	rawRepository RawRepository
}

func NewUsecase(repository *db.Queries, rawRepository RawRepository) Usecase {
	return &usecase{repository: repository, rawRepository: rawRepository}
}

func (u usecase) Store(ctx context.Context, req Request) error {
	return u.repository.SportCreate(ctx, db.SportCreateParams{
		Name:        req.Name,
		Description: req.Description,
		Type:        req.SportType,
		Thumbnail:   req.Thumbnail,
	})
}

func (u usecase) FetchAll(ctx context.Context, page int32, pageSize int32, queries FetchAllQueryParams) (tools.Pagination, error) {
	skip := tools.PaginationSkip(page, pageSize)
	sports, err := u.rawRepository.SportFetchAll(ctx, fetchAllParams{
		Limit:    pageSize,
		Offset:   skip,
		Name:     "%" + queries.Keyword + "%",
		Category: string(queries.Category),
	})
	if err != nil {
		return tools.Pagination{}, fmt.Errorf("error in fetch all sports : %s", err.Error())
	}
	count, err := u.rawRepository.SportCountAll(ctx, countAllParams{
		Name:     "%" + queries.Keyword + "%",
		Category: string(queries.Category),
	})
	if err != nil {
		return tools.Pagination{}, fmt.Errorf("error in count all sports : %s", err.Error())
	}
	return tools.Pagination{
		TotalItem: count,
		PageSize:  pageSize,
		Page:      page,
		Data:      sports,
	}, nil
}

func (u usecase) Update(ctx context.Context, req Request, sportID string) error {
	sportUUID, err := uuid.Parse(sportID)
	if err != nil {
		return fmt.Errorf("error parsing sport id : %s", err.Error())
	}
	if _, err = u.repository.SportCheckOne(ctx, sportUUID); err != nil {
		return fmt.Errorf("error in check sport : %s", err.Error())
	}
	return u.repository.SportUpdate(ctx, db.SportUpdateParams{
		Name:        req.Name,
		Description: req.Description,
		Type:        req.SportType,
		Thumbnail:   req.Thumbnail,
		ID:          sportUUID,
	})
}

func (u usecase) Delete(ctx context.Context, sportID string) error {
	sportUUID, err := uuid.Parse(sportID)
	if err != nil {
		return fmt.Errorf("error parsing sport id : %s", err.Error())
	}
	if _, err = u.repository.SportCheckOne(ctx, sportUUID); err != nil {
		return fmt.Errorf("error in check sport : %s", err.Error())
	}
	return u.repository.SportDelete(ctx, sportUUID)
}
