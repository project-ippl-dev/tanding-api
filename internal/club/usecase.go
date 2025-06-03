package club

//go:generate mockgen -source=./usecase.go -destination=../../mocks/club/usecase_mock.go

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
	"net/http"

	"github.com/google/uuid"

	"github.com/project-ippl-dev/tanding-api/internal/db"
)

type Usecase interface {
	Store(ctx context.Context, req Request, userID string) (uuid.UUID, error)
	Update(ctx context.Context, req Request, userID string, clubID uuid.UUID) error
	Delete(ctx context.Context, userID string, clubID uuid.UUID) error
	Invite(ctx context.Context, req ParticipantReq, userID string) error
	Join(ctx context.Context, req JoinParam, userID string) (statusCode int, err error)
	FetchJoinApproval(ctx context.Context, limit int32, clubID uuid.UUID, ID int64) ([]db.ClubParticipantFetchJoinApprovalRow, error)
	FetchInviteApproval(ctx context.Context, limit int32, userID string, ID int64) ([]db.ClubParticipantFetchInviteApprovalRow, error)
	UpdateJoinApproval(ctx context.Context, userID string, req UpdateJoinApprovalArgs) (statusCode int, err error)
	UpdateInviteApproval(ctx context.Context, userID string, req UpdateInviteApprovalArgs) error
	FetchParticipant(ctx context.Context, req FetchParticipantParam) (tools.Pagination, error)
	Handover(ctx context.Context, userID string, req HandoverReq) (statusCode int, err error)
	FetchAll(ctx context.Context, page int32, pageSize int32, sportID string) (tools.Pagination, error)
	FetchOne(ctx context.Context, clubID uuid.UUID, userID string) (FetchOneResponse, error)
	FetchOwner(ctx context.Context, userID string) ([]db.ClubFetchAllOwnerRow, error)
	Kick(ctx context.Context, arg KickParams, userID string) (statusCode int, err error)
}

type usecase struct {
	repository    *db.Queries
	rawRepository RawRepository
}

func NewUsecase(repository *db.Queries, rawRepository RawRepository) Usecase {
	return &usecase{repository: repository, rawRepository: rawRepository}
}

func (u usecase) Store(ctx context.Context, req Request, userID string) (uuid.UUID, error) {
	tx, err := u.rawRepository.db.Begin()
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("error in start db transaction")
	}
	txQuery := u.repository.WithTx(tx)
	clubID, err := txQuery.ClubCreate(ctx, db.ClubCreateParams{
		Name:      req.Name,
		Logo:      req.Logo,
		Phone:     req.Phone,
		ShortName: req.ShortName,
		UserID:    uuid.MustParse(userID),
	})
	if err != nil {
		if err = tx.Rollback(); err != nil {
			return uuid.UUID{}, fmt.Errorf("error in rollback transaction %s", err.Error())
		}
		return uuid.UUID{}, fmt.Errorf("error in db transaction create club : %s", err.Error())
	}
	for _, sportData := range req.Sports {
		if err = txQuery.ClubAttachSport(ctx, db.ClubAttachSportParams{
			ClubID:  clubID,
			SportID: uuid.MustParse(sportData.SportID),
		}); err != nil {
			if err = tx.Rollback(); err != nil {
				return uuid.UUID{}, fmt.Errorf("error in rollback transaction %s", err.Error())
			}
			return uuid.UUID{}, fmt.Errorf("error in db transaction attach sport to club : %s", err.Error())
		}
	}
	if err = txQuery.ClubParticipantCreate(ctx, db.ClubParticipantCreateParams{
		ClubID: clubID,
		UserID: uuid.MustParse(userID),
	}); err != nil {
		if err = tx.Rollback(); err != nil {
			return uuid.UUID{}, fmt.Errorf("error in rollback transaction %s", err.Error())
		}
		return uuid.UUID{}, fmt.Errorf("error in db transaction add participant to club : %s", err.Error())
	}
	if err = tx.Commit(); err != nil {
		return uuid.UUID{}, fmt.Errorf("error in commit transaction")
	}
	return clubID, nil
}

func (u usecase) Update(ctx context.Context, req Request, userID string, clubID uuid.UUID) error {
	ID, err := u.repository.ClubCheckOne(ctx, db.ClubCheckOneParams{
		ID:     clubID,
		UserID: uuid.MustParse(userID),
	})
	if err != nil {
		return fmt.Errorf("error in check club : %s", err.Error())
	}
	return u.repository.ClubUpdate(ctx, db.ClubUpdateParams{
		Name:      req.Name,
		Logo:      req.Logo,
		Phone:     req.Phone,
		ShortName: req.ShortName,
		ID:        ID,
	})
}

func (u usecase) Delete(ctx context.Context, userID string, clubID uuid.UUID) error {
	ID, err := u.repository.ClubCheckOne(ctx, db.ClubCheckOneParams{
		ID:     clubID,
		UserID: uuid.MustParse(userID),
	})
	if err != nil {
		return fmt.Errorf("error in check club : %s", err.Error())
	}
	return u.repository.ClubDelete(ctx, ID)
}

func (u usecase) Invite(ctx context.Context, req ParticipantReq, userID string) error {
	ID, err := u.repository.ClubCheckOne(ctx, db.ClubCheckOneParams{
		ID:     uuid.MustParse(req.ClubID),
		UserID: uuid.MustParse(userID),
	})
	if err != nil {
		return fmt.Errorf("error in check club : %s", err.Error())
	}
	tx, err := u.rawRepository.db.Begin()
	if err != nil {
		return fmt.Errorf("error in start db transaction : %s", err.Error())
	}
	txQuery := u.repository.WithTx(tx)
	for _, participant := range req.Participants {
		if _, err = u.repository.ClubParticipantCheckParticipant(ctx, db.ClubParticipantCheckParticipantParams{
			ClubID: uuid.MustParse(req.ClubID),
			UserID: uuid.MustParse(participant.SportID),
		}); err == nil {
			return fmt.Errorf("you already join club")
		}
		if err = txQuery.ClubParticipantInvite(ctx, db.ClubParticipantInviteParams{
			ClubID:  ID,
			UserID:  uuid.MustParse(participant.UserID),
			SportID: uuid.MustParse(participant.SportID),
		}); err != nil {
			if err = tx.Rollback(); err != nil {
				return fmt.Errorf("error in rollback transaction %s", err.Error())
			}
			return fmt.Errorf("error in club participant invite: %s", err.Error())
		}
	}
	return tx.Commit()
}

func (u usecase) Join(ctx context.Context, req JoinParam, userID string) (statusCode int, err error) {
	if _, err = u.repository.ClubCheckOneWithoutUserID(ctx, uuid.MustParse(req.ClubID)); err != nil {
		return http.StatusNotFound, fmt.Errorf("error in check club : %s", err.Error())
	}
	if _, err = u.repository.ClubParticipantCheckParticipant(ctx, db.ClubParticipantCheckParticipantParams{
		ClubID: uuid.MustParse(req.ClubID),
		UserID: uuid.MustParse(userID),
	}); err == nil {
		return http.StatusBadRequest, fmt.Errorf("you already join club")
	}
	if err = u.repository.ClubParticipantJoin(ctx, db.ClubParticipantJoinParams{
		ClubID:  uuid.MustParse(req.ClubID),
		UserID:  uuid.MustParse(userID),
		SportID: uuid.MustParse(req.SportID),
	}); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error in store club participant join : %s", err.Error())
	}
	return http.StatusOK, nil
}

func (u usecase) FetchJoinApproval(ctx context.Context, limit int32, clubID uuid.UUID, ID int64) ([]db.ClubParticipantFetchJoinApprovalRow, error) {
	if ID == 0 {
		LatestID, _ := u.repository.ClubParticipantFetchLatestIDJoinApproval(ctx, clubID)
		ID = LatestID + 1
	}
	data, err := u.repository.ClubParticipantFetchJoinApproval(ctx, db.ClubParticipantFetchJoinApprovalParams{
		ClubID: clubID,
		ID:     ID,
		Limit:  limit,
	})
	if err != nil {
		return nil, fmt.Errorf("error in fetch join approval : %s", err.Error())
	}
	return data, nil
}

func (u usecase) FetchInviteApproval(ctx context.Context, limit int32, userID string, ID int64) ([]db.ClubParticipantFetchInviteApprovalRow, error) {
	if ID == 0 {
		LatestID, _ := u.repository.ClubParticipantFetchLatestIDInviteApproval(ctx, uuid.MustParse(userID))
		ID = LatestID + 1
	}
	data, err := u.repository.ClubParticipantFetchInviteApproval(ctx, db.ClubParticipantFetchInviteApprovalParams{
		UserID: uuid.MustParse(userID),
		ID:     ID,
		Limit:  limit,
	})
	if err != nil {
		return nil, fmt.Errorf("error in fetch invite approval : %s", err.Error())
	}
	return data, nil
}

func (u usecase) UpdateJoinApproval(ctx context.Context, userID string, req UpdateJoinApprovalArgs) (statusCode int, err error) {
	clubID, err := u.repository.ClubCheckOne(ctx, db.ClubCheckOneParams{
		ID:     uuid.MustParse(req.ClubID),
		UserID: uuid.MustParse(userID),
	})
	if err != nil {
		return http.StatusNotFound, fmt.Errorf("club not found : %s", err.Error())
	}
	if _, err = u.repository.ClubParticipantCheckInviteApproval(ctx, db.ClubParticipantCheckInviteApprovalParams{
		ID:     req.ApprovalID,
		UserID: uuid.MustParse(userID),
	}); err != nil {
		return http.StatusNotFound, fmt.Errorf("club participant join not found : %s", err.Error())
	}
	if err = u.repository.ClubParticipantUpdateJoinApproval(ctx, db.ClubParticipantUpdateJoinApprovalParams{
		ClubApproval: sql.NullBool{
			Bool:  *req.Status,
			Valid: true,
		},
		ID:     req.ApprovalID,
		ClubID: clubID,
	}); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error in update join approval : %s", err.Error())
	}
	return http.StatusOK, nil
}

func (u usecase) UpdateInviteApproval(ctx context.Context, userID string, req UpdateInviteApprovalArgs) error {
	if _, err := u.repository.ClubParticipantCheckJoinApproval(ctx, db.ClubParticipantCheckJoinApprovalParams{
		ID:     req.ApprovalID,
		UserID: uuid.MustParse(userID),
	}); err != nil {
		return fmt.Errorf("error in check participant invite approval : %s", err.Error())
	}
	return u.repository.ClubParticipantUpdateInviteApproval(ctx, db.ClubParticipantUpdateInviteApprovalParams{
		UserApproval: sql.NullBool{
			Bool:  *req.Status,
			Valid: true,
		},
		ID:     req.ApprovalID,
		UserID: uuid.MustParse(userID),
	})
}

func (u usecase) FetchParticipant(ctx context.Context, req FetchParticipantParam) (tools.Pagination, error) {
	skip := tools.PaginationSkip(req.Page, req.PageSize)
	queryParams := participantQueryParams{
		SportID:        req.SportID,
		CanParticipate: req.CanParticipate,
		ClubID:         req.ClubID,
	}

	totalPoint, _ := u.repository.RankFetchPointByClubID(ctx, req.ClubID)

	count, err := u.rawRepository.ClubParticipantCountAll(ctx, queryParams)
	if err != nil {
		return tools.Pagination{}, fmt.Errorf("error in count club participants : %s", err.Error())
	}

	participants, err := u.rawRepository.ClubParticipantFetchAll(ctx, participantFetchAllParams{
		participantQueryParams: queryParams,
		Limit:                  req.PageSize,
		Offset:                 skip,
	})
	if err != nil {
		return tools.Pagination{}, fmt.Errorf("error in fetch club participants : %s", err.Error())
	}

	result := []FetchParticipantRow{}
	for _, participant := range participants {
		point, _ := u.repository.RankFetchPointByClubIDAndUserID(ctx, db.RankFetchPointByClubIDAndUserIDParams{
			ClubID: req.ClubID,
			UserID: participant.UserID,
		})
		result = append(result, FetchParticipantRow{
			ParticipantFetchAllRow: participant,
			Point:                  point,
		})
	}

	return tools.Pagination{
		TotalItem: count,
		PageSize:  req.PageSize,
		Page:      req.Page,
		Data: FetchParticipantResponse{
			TotalPoint:   totalPoint,
			Participants: result,
		},
	}, nil
}

func (u usecase) Handover(ctx context.Context, userID string, req HandoverReq) (statusCode int, err error) {
	clubID, err := u.repository.ClubCheckOne(ctx, db.ClubCheckOneParams{
		ID:     uuid.MustParse(req.ClubID),
		UserID: uuid.MustParse(userID),
	})
	if err != nil {
		return http.StatusNotFound, fmt.Errorf("error in check club : %s", err.Error())
	}
	if _, err := u.repository.ClubParticipantCheckOne(ctx, db.ClubParticipantCheckOneParams{
		ClubID: uuid.MustParse(req.ClubID),
		UserID: uuid.MustParse(req.UserID),
	}); err != nil {
		return http.StatusNotFound, fmt.Errorf("error in club participant check one : %s", err.Error())
	}
	if err := u.repository.ClubHandover(ctx, db.ClubHandoverParams{
		UserID: uuid.MustParse(req.UserID),
		ID:     clubID,
	}); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error in update club handover : %s", err.Error())
	}
	return http.StatusOK, nil
}

func (u usecase) FetchAll(ctx context.Context, page int32, pageSize int32, sportID string) (tools.Pagination, error) {
	skip := tools.PaginationSkip(page, pageSize)
	count, err := u.rawRepository.ClubCountAll(ctx, sportID)
	data := []FetchAllResponse{}
	if err != nil {
		return tools.Pagination{}, fmt.Errorf("error in count all club : %s", err.Error())
	}
	clubs, err := u.rawRepository.ClubFetchAll(ctx, fetchAllParams{
		Limit:   pageSize,
		Offset:  skip,
		SportID: sportID,
	})
	if err != nil {
		return tools.Pagination{}, fmt.Errorf("error in fetch all club : %s", err.Error())
	}

	for _, club := range clubs {
		sports, err := u.repository.ClubSportFetchByClubID(ctx, club.ID)
		if err != nil {
			return tools.Pagination{}, fmt.Errorf("error in fetch sports in club : %s", err.Error())
		}
		clubResponse := FetchAllResponse{
			FetchAllRow: club,
			Sports:      sports,
		}
		data = append(data, clubResponse)
	}

	return tools.Pagination{
		TotalItem: count,
		PageSize:  pageSize,
		Page:      page,
		Data:      data,
	}, nil
}

func (u usecase) FetchOne(ctx context.Context, clubID uuid.UUID, userID string) (FetchOneResponse, error) {
	club, err := u.repository.ClubFetchOne(ctx, clubID)
	if err != nil {
		return FetchOneResponse{}, fmt.Errorf("error in fetch one club : %s", err.Error())
	}
	sports, err := u.repository.ClubSportFetchByClubID(ctx, clubID)
	if err != nil {
		return FetchOneResponse{}, fmt.Errorf("error in fetch club sports : %s", err.Error())
	}

	var privilege, joined bool
	if _, err := u.repository.ClubFetchOneOwner(ctx, db.ClubFetchOneOwnerParams{
		ID:     club.ID,
		UserID: uuid.MustParse(userID),
	}); err == nil {
		privilege = true
	}

	if _, err := u.repository.ClubParticipantCheckParticipant(ctx, db.ClubParticipantCheckParticipantParams{
		ClubID: clubID,
		UserID: uuid.MustParse(userID),
	}); err == nil {
		joined = true
	}

	return FetchOneResponse{
		ClubFetchOneRow: club,
		Sports:          sports,
		Privilege:       privilege,
		Joined:          joined,
	}, nil
}

func (u usecase) FetchOwner(ctx context.Context, userID string) ([]db.ClubFetchAllOwnerRow, error) {
	return u.repository.ClubFetchAllOwner(ctx, uuid.MustParse(userID))
}

func (u usecase) Kick(ctx context.Context, arg KickParams, userID string) (statusCode int, err error) {
	clubID, err := u.repository.ClubCheckOne(ctx, db.ClubCheckOneParams{
		ID:     arg.ClubID,
		UserID: uuid.MustParse(userID),
	})
	if err != nil {
		return http.StatusForbidden, fmt.Errorf("only owner of the club can kick")
	}
	clubParticipantID, err := u.repository.ClubParticipantCheckByOwner(ctx, db.ClubParticipantCheckByOwnerParams{
		ClubID: clubID,
		UserID: uuid.MustParse(userID),
		ID:     arg.ClubParticipantID,
	})
	if err != nil {
		return http.StatusNotFound, fmt.Errorf("member not found : %s", err.Error())
	}
	if err := u.repository.ClubParticipantKick(ctx, db.ClubParticipantKickParams{
		ID:     clubParticipantID,
		ClubID: clubID,
	}); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error in kick club member : %s", err.Error())
	}
	return http.StatusOK, nil
}
