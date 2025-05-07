package sport

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/project-ippl-dev/tanding-api/internal/db"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
)

type Usecase struct {
	repository    *db.Queries
	rawRepository RawRepository
}

func NewUsecase(repository *db.Queries, rawRepository RawRepository) Usecase {
	return Usecase{repository: repository, rawRepository: rawRepository}
}

func (u Usecase) store(ctx context.Context, req request) error {
	return u.repository.SportCreate(ctx, db.SportCreateParams{
		Name:        req.Name,
		Description: req.Description,
		Type:        req.SportType,
		Thumbnail:   req.Thumbnail,
	})
}

func (u Usecase) fetchAll(ctx context.Context, page int32, pageSize int32, queries fetchAllQueryParams) (tools.Pagination, error) {
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

func (u Usecase) update(ctx context.Context, req request, sportID string) error {
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

func (u Usecase) delete(ctx context.Context, sportID string) error {
	sportUUID, err := uuid.Parse(sportID)
	if err != nil {
		return fmt.Errorf("error parsing sport id : %s", err.Error())
	}
	if _, err = u.repository.SportCheckOne(ctx, sportUUID); err != nil {
		return fmt.Errorf("error in check sport : %s", err.Error())
	}
	return u.repository.SportDelete(ctx, sportUUID)
}
