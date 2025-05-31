package certificate

//go:generate mockgen -source=./usecase.go -destination=../../mocks/certificate/usecase_mock.go

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/project-ippl-dev/tanding-api/internal/db"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
)

type Usecase interface {
	FetchOne(ctx context.Context, certificateID uuid.UUID) (Response, error)
	FetchByUserID(ctx context.Context, page, pageSize int32, userID string) (tools.Pagination, error)
	FetchByClubID(ctx context.Context, page, pageSize int32, clubID uuid.UUID) (tools.Pagination, error)
}
type usecase struct {
	repository *db.Queries
}

func NewUsecase(repository *db.Queries) Usecase {
	return &usecase{repository: repository}
}

func (u usecase) FetchOne(ctx context.Context, certificateID uuid.UUID) (Response, error) {
	result, err := u.repository.CertificateFetchOne(ctx, certificateID)
	if err != nil {
		return Response{}, fmt.Errorf("certificate not found : %s", err.Error())
	}
	event, err := u.repository.EventFetchOneInfiniteByID(ctx, result.EventID)
	if err != nil {
		return Response{}, fmt.Errorf("event not found : %s", err.Error())
	}
	participants, err := u.repository.EventRegistrationCountAllByStatusApproved(ctx, event.ID)
	if err != nil {
		return Response{}, fmt.Errorf("event not found : %s", err.Error())
	}
	photo, err := u.repository.UserFetchPhotoByID(ctx, result.UserID)
	if err != nil {
		return Response{}, fmt.Errorf("recipient not found : %s", err.Error())
	}
	return Response{
		Certificate: result,
		Event: EventDetail{
			EventFetchOneInfiniteByIDRow: event,
			Participants:                 participants,
		},
		Recipient: photo,
	}, nil
}

func (u usecase) FetchByUserID(ctx context.Context, page, pageSize int32, userID string) (tools.Pagination, error) {
	skip := tools.PaginationSkip(page, pageSize)

	count, err := u.repository.CertificateCountAllByUserID(ctx, uuid.MustParse(userID))
	if err != nil {
		return tools.Pagination{}, fmt.Errorf("error in count certificate : %s", err.Error())
	}
	certificates, _ := u.repository.CertificateFetchAllByUserID(ctx, db.CertificateFetchAllByUserIDParams{
		UserID: uuid.MustParse(userID),
		Limit:  pageSize,
		Offset: skip,
	})

	return tools.Pagination{
		TotalItem: count,
		PageSize:  pageSize,
		Page:      page,
		Data:      certificates,
	}, nil
}

func (u usecase) FetchByClubID(ctx context.Context, page, pageSize int32, clubID uuid.UUID) (tools.Pagination, error) {
	skip := tools.PaginationSkip(page, pageSize)

	count, err := u.repository.ClubCertificateCountAllByClubID(ctx, clubID)
	if err != nil {
		return tools.Pagination{}, fmt.Errorf("error in count certificate : %s", err.Error())
	}
	certificates, _ := u.repository.ClubCertificateFetchAllByUserID(ctx, db.ClubCertificateFetchAllByUserIDParams{
		ClubID: clubID,
		Limit:  pageSize,
		Offset: skip,
	})

	return tools.Pagination{
		TotalItem: count,
		PageSize:  pageSize,
		Page:      page,
		Data:      certificates,
	}, nil
}
