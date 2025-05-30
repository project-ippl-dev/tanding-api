package classCompetitionRule

//go:generate mockgen -source=./usecase.go -destination=../../mocks/classCompetitionRule/usecase_mock.go

import (
	"context"
	"fmt"

	"github.com/project-ippl-dev/tanding-api/internal/db"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
)

type Usecase interface {
	Store(ctx context.Context, req Request) error
	FetchAll(ctx context.Context, page int32, pageSize int32) (tools.Pagination, error)
	FetchOne(ctx context.Context, ruleID int64) (db.ClassRuleFetchOneRow, error)
	Update(ctx context.Context, req Request, ruleID int64) error
	Delete(ctx context.Context, ruleID int64) error
}

type usecase struct {
	repository *db.Queries
}

func NewUsecase(repository *db.Queries) Usecase {
	return &usecase{repository: repository}
}

func (u usecase) Store(ctx context.Context, req Request) error {
	return u.repository.ClassRuleCreate(ctx, db.ClassRuleCreateParams{
		Name:   req.Name,
		Male:   req.Male,
		Female: req.Female,
		Total:  req.Total,
	})
}

func (u usecase) FetchAll(ctx context.Context, page int32, pageSize int32) (tools.Pagination, error) {
	skip := tools.PaginationSkip(page, pageSize)
	rules, err := u.repository.ClassRuleFetchAll(ctx, db.ClassRuleFetchAllParams{
		Limit:  pageSize,
		Offset: skip,
	})
	if err != nil {
		return tools.Pagination{}, fmt.Errorf("error in fetch all competition rules : %s", err.Error())
	}
	count, err := u.repository.ClassRuleCountAll(ctx)
	if err != nil {
		return tools.Pagination{}, fmt.Errorf("error in count all competition rules : %s", err.Error())
	}
	return tools.Pagination{
		TotalItem: count,
		PageSize:  pageSize,
		Page:      page,
		Data:      rules,
	}, nil
}

func (u usecase) FetchOne(ctx context.Context, ruleID int64) (db.ClassRuleFetchOneRow, error) {
	return u.repository.ClassRuleFetchOne(ctx, ruleID)
}

func (u usecase) Update(ctx context.Context, req Request, ruleID int64) error {
	if _, err := u.repository.ClassRuleFetchOne(ctx, ruleID); err != nil {
		return fmt.Errorf("error in check class rules : %s", err.Error())
	}
	return u.repository.ClassRuleUpdate(ctx, db.ClassRuleUpdateParams{
		Name:   req.Name,
		Male:   req.Male,
		Female: req.Female,
		Total:  req.Total,
		ID:     ruleID,
	})
}

func (u usecase) Delete(ctx context.Context, ruleID int64) error {
	if _, err := u.repository.ClassRuleFetchOne(ctx, ruleID); err != nil {
		return fmt.Errorf("error in check class rules : %s", err.Error())
	}
	return u.repository.ClassRuleDelete(ctx, ruleID)
}
