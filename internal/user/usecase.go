package user

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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

func (u Usecase) search(ctx context.Context, args searchParams, userID string) ([]db.UserFetchByKeywordRow, error) {
	return u.repository.UserFetchByKeyword(ctx, db.UserFetchByKeywordParams{
		Name:     "%" + args.Keyword + "%",
		Username: "%" + args.Keyword + "%",
		Limit:    args.Limit,
		ID:       uuid.MustParse(userID),
	})
}

func (u Usecase) fetchOne(ctx context.Context, userID string) (basicInformationResponse, error) {
	basic, err := u.repository.UserFetchBasicInformation(ctx, uuid.MustParse(userID))
	if err != nil {
		return basicInformationResponse{}, fmt.Errorf("error in fetch basic information : " + err.Error())
	}
	return basicInformationResponse{
		UserFetchBasicInformationRow: basic,
	}, nil
}

func (u Usecase) update(ctx context.Context, arg updateBasicInformationParams, userID string) error {
	bornOn, _ := time.Parse("2006-01-02", arg.BornOn)
	var bornOnStatus bool
	if !bornOn.IsZero() {
		bornOnStatus = true
	}

	if err := u.repository.UserUpdateBasicInformation(ctx, db.UserUpdateBasicInformationParams{
		Name:   arg.Name,
		BornAt: arg.BornAt,
		BornOn: sql.NullTime{
			Time:  bornOn,
			Valid: bornOnStatus,
		},
		IdentityNumber: arg.IdentityNumber,
		Phone:          arg.Phone,
		Photo:          arg.Photo,
		Gender:         arg.Gender,
		About:          arg.About,
		ID:             uuid.MustParse(userID),
	}); err != nil {
		return fmt.Errorf("error in update basic information : " + err.Error())
	}

	return nil
}

func (u Usecase) fetchAll(ctx context.Context, page, pageSize int32) (tools.Pagination, error) {
	skip := tools.PaginationSkip(page, pageSize)
	count, err := u.repository.UserCountAll(ctx)
	if err != nil {
		return tools.Pagination{}, fmt.Errorf("error in count user : " + err.Error())
	}
	users, err := u.repository.UserFetchAll(ctx, db.UserFetchAllParams{
		Limit:  pageSize,
		Offset: skip,
	})
	if err != nil {
		return tools.Pagination{}, fmt.Errorf("error in fetch user : " + err.Error())
	}

	return tools.Pagination{
		TotalItem: count,
		PageSize:  pageSize,
		Page:      page,
		Data:      users,
	}, nil
}

func (u Usecase) fetchLastLogin(ctx context.Context, page, pageSize int32) (tools.Pagination, error) {
	skip := tools.PaginationSkip(page, pageSize)
	count, err := u.repository.LoginDetailCount(ctx)
	if err != nil {
		return tools.Pagination{}, fmt.Errorf("error count last login : " + err.Error())
	}
	lastLogin, err := u.repository.LoginDetailFetchAll(ctx, db.LoginDetailFetchAllParams{
		Limit:  pageSize,
		Offset: skip,
	})
	if err != nil {
		return tools.Pagination{}, fmt.Errorf("error fetch last login : " + err.Error())
	}
	return tools.Pagination{
		TotalItem: count,
		PageSize:  pageSize,
		Page:      page,
		Data:      lastLogin,
	}, nil
}
