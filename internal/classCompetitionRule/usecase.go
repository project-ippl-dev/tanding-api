package classCompetitionRule

import (
	"context"
	"fmt"

	"github.com/project-ippl-dev/tanding-api/internal/db"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
)

type Usecase struct {
	repository *db.Queries
}

func NewUsecase(repository *db.Queries) Usecase {
	return Usecase{repository: repository}
}

func (u Usecase) store(ctx context.Context, req request) error {
	return u.repository.ClassRuleCreate(ctx, db.ClassRuleCreateParams{
		Name:   req.Name,
		Male:   req.Male,
		Female: req.Female,
		Total:  req.Total,
	})
}

func (u Usecase) fetchAll(ctx context.Context, page int32, pageSize int32) (tools.Pagination, error) {
	skip := tools.PaginationSkip(page, pageSize)
	rules, err := u.repository.ClassRuleFetchAll(ctx, db.ClassRuleFetchAllParams{
		Limit:  pageSize,
		Offset: skip,
	})
	if err != nil {
		return tools.Pagination{}, fmt.Errorf("error in fetch all competition rules : " + err.Error())
	}
	count, err := u.repository.ClassRuleCountAll(ctx)
	if err != nil {
		return tools.Pagination{}, fmt.Errorf("error in count all competition rules : " + err.Error())
	}
	return tools.Pagination{
		TotalItem: count,
		PageSize:  pageSize,
		Page:      page,
		Data:      rules,
	}, nil
}

func (u Usecase) fetchOne(ctx context.Context, ruleID int64) (db.ClassRuleFetchOneRow, error) {
	return u.repository.ClassRuleFetchOne(ctx, ruleID)
}

func (u Usecase) update(ctx context.Context, req request, ruleID int64) error {
	if _, err := u.repository.ClassRuleFetchOne(ctx, ruleID); err != nil {
		return fmt.Errorf("error in check class rules : " + err.Error())
	}
	return u.repository.ClassRuleUpdate(ctx, db.ClassRuleUpdateParams{
		Name:   req.Name,
		Male:   req.Male,
		Female: req.Female,
		Total:  req.Total,
		ID:     ruleID,
	})
}

func (u Usecase) delete(ctx context.Context, ruleID int64) error {
	if _, err := u.repository.ClassRuleFetchOne(ctx, ruleID); err != nil {
		return fmt.Errorf("error in check class rules : " + err.Error())
	}
	return u.repository.ClassRuleDelete(ctx, ruleID)
}
