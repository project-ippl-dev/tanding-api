package user

//go:generate mockgen -source=./usecase.go -destination=../../mocks/user/usecase_mock.go

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/project-ippl-dev/tanding-api/internal/db"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
)

type Usecase interface {
	Search(ctx context.Context, args SearchParams, userID string) ([]db.UserFetchByKeywordRow, error)
	FetchOne(ctx context.Context, userID string) (BasicInformationResponse, error)
	Update(ctx context.Context, arg UpdateBasicInformationParams, userID string) error
	FetchAll(ctx context.Context, page, pageSize int32) (tools.Pagination, error)
	FetchLastLogin(ctx context.Context, page, pageSize int32) (tools.Pagination, error)
}

type usecase struct {
	repository    *db.Queries
	rawRepository RawRepository
}

func NewUsecase(repository *db.Queries, rawRepository RawRepository) Usecase {
	return &usecase{repository: repository, rawRepository: rawRepository}
}

func (u usecase) Search(ctx context.Context, args SearchParams, userID string) ([]db.UserFetchByKeywordRow, error) {
	return u.repository.UserFetchByKeyword(ctx, db.UserFetchByKeywordParams{
		Name:     "%" + args.Keyword + "%",
		Username: "%" + args.Keyword + "%",
		Limit:    args.Limit,
		ID:       uuid.MustParse(userID),
	})
}

func (u usecase) FetchOne(ctx context.Context, userID string) (BasicInformationResponse, error) {
	basic, err := u.repository.UserFetchBasicInformation(ctx, uuid.MustParse(userID))
	if err != nil {
		return BasicInformationResponse{}, fmt.Errorf("error in fetch basic information : %s", err.Error())
	}
	return BasicInformationResponse{
		UserFetchBasicInformationRow: basic,
	}, nil
}

func (u usecase) Update(ctx context.Context, arg UpdateBasicInformationParams, userID string) error {
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
		return fmt.Errorf("error in update basic information : %s", err.Error())
	}

	return nil
}

func (u usecase) FetchAll(ctx context.Context, page, pageSize int32) (tools.Pagination, error) {
	skip := tools.PaginationSkip(page, pageSize)
	count, err := u.repository.UserCountAll(ctx)
	if err != nil {
		return tools.Pagination{}, fmt.Errorf("error in count user : %s", err.Error())
	}
	users, err := u.repository.UserFetchAll(ctx, db.UserFetchAllParams{
		Limit:  pageSize,
		Offset: skip,
	})
	if err != nil {
		return tools.Pagination{}, fmt.Errorf("error in fetch user : %s", err.Error())
	}

	return tools.Pagination{
		TotalItem: count,
		PageSize:  pageSize,
		Page:      page,
		Data:      users,
	}, nil
}

func (u usecase) FetchLastLogin(ctx context.Context, page, pageSize int32) (tools.Pagination, error) {
	skip := tools.PaginationSkip(page, pageSize)
	count, err := u.repository.LoginDetailCount(ctx)
	if err != nil {
		return tools.Pagination{}, fmt.Errorf("error count last login : %s", err.Error())
	}
	lastLogin, err := u.repository.LoginDetailFetchAll(ctx, db.LoginDetailFetchAllParams{
		Limit:  pageSize,
		Offset: skip,
	})
	if err != nil {
		return tools.Pagination{}, fmt.Errorf("error fetch last login : %s", err.Error())
	}
	return tools.Pagination{
		TotalItem: count,
		PageSize:  pageSize,
		Page:      page,
		Data:      lastLogin,
	}, nil
}
