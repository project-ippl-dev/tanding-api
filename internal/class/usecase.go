package class

//go:generate mockgen -source=./usecase.go -destination=../../mocks/class/usecase_mock.go

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/project-ippl-dev/tanding-api/internal/db"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
)

type Usecase interface {
	Store(ctx context.Context, req Request) error
	FetchAll(ctx context.Context, page int32, pageSize int32) (tools.Pagination, error)
	FetchBySportID(ctx context.Context, page int32, pageSize int32, sportID uuid.UUID) (tools.Pagination, error)
	Update(ctx context.Context, req Request, classID uuid.UUID) error
	Delete(ctx context.Context, decoded tools.JWT, classID uuid.UUID) (statusCode int, err error)
}

type usecase struct {
	repository *db.Queries
}

func NewUsecase(repository *db.Queries) Usecase {
	return &usecase{repository: repository}
}

func (u usecase) Store(ctx context.Context, req Request) error {
	return u.repository.ClassCreate(ctx, db.ClassCreateParams{
		SportID:                uuid.MustParse(req.SportID),
		Name:                   req.Name,
		ClassCompetitionRuleID: req.ClassRuleID,
		Type:                   req.Type,
		MatchType:              req.MatchType,
	})
}

func (u usecase) FetchAll(ctx context.Context, page int32, pageSize int32) (tools.Pagination, error) {
	skip := tools.PaginationSkip(page, pageSize)
	classes, err := u.repository.ClassFetchAll(ctx, db.ClassFetchAllParams{
		Limit:  pageSize,
		Offset: skip,
	})
	if err != nil {
		return tools.Pagination{}, fmt.Errorf("error in fetch class : %s", err.Error())
	}
	count, err := u.repository.ClassCountAll(ctx)
	if err != nil {
		return tools.Pagination{}, fmt.Errorf("error in count class : %s", err.Error())
	}
	return tools.Pagination{
		TotalItem: count,
		PageSize:  pageSize,
		Page:      page,
		Data:      classes,
	}, nil
}

func (u usecase) FetchBySportID(ctx context.Context, page int32, pageSize int32, sportID uuid.UUID) (tools.Pagination, error) {
	skip := tools.PaginationSkip(page, pageSize)
	classes, err := u.repository.ClassFetchBySportID(ctx, db.ClassFetchBySportIDParams{
		SportID: sportID,
		Limit:   pageSize,
		Offset:  skip,
	})
	if err != nil {
		return tools.Pagination{}, fmt.Errorf("error in fetch class : %s", err.Error())
	}
	count, err := u.repository.ClassCountBySportID(ctx, sportID)
	if err != nil {
		return tools.Pagination{}, fmt.Errorf("error in fetch class : %s", err.Error())
	}
	return tools.Pagination{
		TotalItem: count,
		PageSize:  pageSize,
		Page:      page,
		Data:      classes,
	}, nil
}

func (u usecase) Update(ctx context.Context, req Request, classID uuid.UUID) error {
	if _, err := u.repository.ClassCheckOne(ctx, classID); err != nil {
		return fmt.Errorf("error in check class : %s", err.Error())
	}
	return u.repository.ClassUpdate(ctx, db.ClassUpdateParams{
		SportID:                uuid.MustParse(req.SportID),
		Name:                   req.Name,
		ClassCompetitionRuleID: req.ClassRuleID,
		Type:                   req.Type,
		ID:                     classID,
		MatchType:              req.MatchType,
	})
}

func (u usecase) Delete(ctx context.Context, decoded tools.JWT, classID uuid.UUID) (statusCode int, err error) {
	class, err := u.repository.ClassCheckOne(ctx, classID)
	if err != nil {
		return http.StatusNotFound, fmt.Errorf("error in check class : %s", err.Error())
	}
	if decoded.RoleName == "user" && class.Type != db.ClassTypeCustom {
		return http.StatusForbidden, fmt.Errorf("forbidden access")
	}
	if err = u.repository.ClassDelete(ctx, classID); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error in delete class : %s", err.Error())
	}
	return http.StatusOK, nil
}
