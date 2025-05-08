package eventRegistration

import (
	"context"
	"fmt"
	"net/http"
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
	return Usecase{rawRepository: rawRepository, repository: repository}
}

func (u Usecase) register(ctx context.Context, req registrationRequest) (statusCode int, err error) {
	//Class Rule Validation Based On Gender Event Registration
	classEvent, err := u.repository.ClassEventFetchOne(ctx, uuid.MustParse(req.ClassEventID))
	if err != nil {
		return http.StatusNotFound, fmt.Errorf("class event not found : %s", err.Error())
	}

	status := db.EventRegistrationStatusPending
	if classEvent.Price == 0 {
		status = db.EventRegistrationStatusApproved
	}

	tx, err := u.rawRepository.db.Begin()
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error in start tx : %s", err.Error())
	}
	txQuery := u.repository.WithTx(tx)

	registrationID, err := txQuery.EventRegistrationCreate(ctx, db.EventRegistrationCreateParams{
		ClassEventID:          uuid.MustParse(req.ClassEventID),
		ClubID:                uuid.MustParse(req.ClubID),
		EventID:               uuid.MustParse(req.EventID),
		EventPaymentReceiptID: uuid.UUID{},
		Status:                status,
	})
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return http.StatusInternalServerError, fmt.Errorf("error in rollback tx create event registration : %s", err.Error())
		}
		return http.StatusInternalServerError, fmt.Errorf("error in create event registration : %s", err.Error())
	}
	var userID []string
	for _, member := range req.Members {
		if err := txQuery.EventParticipantCreate(ctx, db.EventParticipantCreateParams{
			EventRegistrationID: registrationID,
			UserID:              member.UserID,
		}); err != nil {
			if err := tx.Rollback(); err != nil {
				return http.StatusInternalServerError, fmt.Errorf("error in rollback tx create event participant : %s", err.Error())
			}
			return http.StatusInternalServerError, fmt.Errorf("error in create event participant : %s", err.Error())
		}
		userID = append(userID, member.UserID.String())
	}

	users, err := u.rawRepository.UserFetchInID(ctx, userID)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return http.StatusInternalServerError, fmt.Errorf("error in rollback tx fetch users : %s", err.Error())
		}
		return http.StatusBadRequest, fmt.Errorf("error in fetch users : %s", err.Error())
	}

	var male, female int16
	for _, user := range users {
		if !user.CanParticipate {
			if err := tx.Rollback(); err != nil {
				return http.StatusInternalServerError, fmt.Errorf("error in rollback tx user validation : %s", err.Error())
			}
			return http.StatusForbidden, fmt.Errorf("please fill your basic information first before continue to register")
		}
		switch user.Gender {
		case "male":
			male += 1
		case "female":
			female += 1
		}
	}

	if classEvent.RuleFemale != 0 {
		if classEvent.RuleFemale != female {
			if err := tx.Rollback(); err != nil {
				return http.StatusInternalServerError, fmt.Errorf("error in rollback tx class rule female validation : %s", err.Error())
			}
			return http.StatusBadRequest, fmt.Errorf("fail class rule validation for female rule")
		}
	}
	if classEvent.RuleMale != 0 {
		if classEvent.RuleMale != male {
			if err := tx.Rollback(); err != nil {
				return http.StatusInternalServerError, fmt.Errorf("error in rollback tx class rule male validation : %s", err.Error())
			}
			return http.StatusBadRequest, fmt.Errorf("fail class rule validation for male rule")
		}
	}

	if classEvent.RuleTotal != male+female {
		if err := tx.Rollback(); err != nil {
			return http.StatusInternalServerError, fmt.Errorf("error in rollback tx class rule male validation : %s", err.Error())
		}
		return http.StatusBadRequest, fmt.Errorf("fail class rule validation for male rule")
	}

	if err := tx.Commit(); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error in commit tx : %s", err.Error())
	}
	return http.StatusCreated, nil
}

func (u Usecase) fetchAll(ctx context.Context, args fetchAllParams, page, pageSize int32) (tools.Pagination, error) {
	skip := tools.PaginationSkip(page, pageSize)

	count, err := u.rawRepository.EventRegistrationCountAll(ctx, fetchQueryParams{args})
	if err != nil {
		return tools.Pagination{}, fmt.Errorf("error in count registration : %s", err.Error())
	}
	registrations, err := u.rawRepository.EventRegistrationFetchAll(ctx, fetchAllDBParams{
		fetchQueryParams: fetchQueryParams{args},
		Limit:            pageSize,
		Offset:           skip,
	})
	if err != nil {
		return tools.Pagination{}, fmt.Errorf("error in fetch registration : %s", err.Error())
	}
	data := []fetchResponse{}
	for _, registration := range registrations {
		participants, err := u.repository.EventParticipantFetchByRegistrationID(ctx, registration.ID)
		if err != nil {
			return tools.Pagination{}, fmt.Errorf("error in fetch participant by specific registration id : %s", err.Error())
		}
		data = append(data, fetchResponse{
			fetchAllRow:  registration,
			Participants: participants,
		})
	}

	return tools.Pagination{
		TotalItem: count,
		PageSize:  pageSize,
		Page:      page,
		Data:      data,
	}, nil
}

func (u Usecase) update(ctx context.Context, req updateRegistrationRequest, userID string) (int, error) {
	registration, err := u.repository.EventRegistrationFetchOne(ctx, req.RegisterID)
	if err != nil {
		return http.StatusNotFound, fmt.Errorf("error in fetch registration : %s", err.Error())
	}

	event, err := u.repository.EventFetchOne(ctx, registration.EventID)
	if err != nil {
		return http.StatusNotFound, fmt.Errorf("error in fetch one registration : %s", err.Error())
	}

	if event.Deadline.Before(time.Now()); err != nil {
		return http.StatusForbidden, fmt.Errorf("can't update a registration after deadline")
	}

	if _, err := u.repository.ClubCheckByIDAndUserID(ctx, db.ClubCheckByIDAndUserIDParams{
		ID:     registration.ClubID,
		UserID: uuid.MustParse(userID),
	}); err != nil {
		return http.StatusForbidden, fmt.Errorf("not auth club owner to update : %s", err.Error())
	}

	if err := u.repository.EventRegistrationUpdate(ctx, db.EventRegistrationUpdateParams{
		ClassEventID: uuid.MustParse(req.ClassEventID),
		ID:           req.RegisterID,
	}); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error in update event registration : %s", err.Error())
	}

	return http.StatusOK, nil
}

func (u Usecase) setReject(ctx context.Context, arg setStatusRequest, userID string) (int, error) {
	registration, err := u.repository.EventRegistrationFetchOne(ctx, arg.RegisterID)
	if err != nil {
		return http.StatusNotFound, fmt.Errorf("error in fetch registration : %s", err.Error())
	}

	if _, err := u.repository.ClubCheckByIDAndUserID(ctx, db.ClubCheckByIDAndUserIDParams{
		ID:     registration.ClubID,
		UserID: uuid.MustParse(userID),
	}); err != nil {
		return http.StatusForbidden, fmt.Errorf("not auth club owner to update : %s", err.Error())
	}

	if err := u.repository.EventRegistrationSetReject(ctx, registration.ID); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error in update status registration to reject : %s", err.Error())
	}

	return http.StatusOK, nil
}

func (u Usecase) fetchParticipant(ctx context.Context, eventID uuid.UUID) ([]fetchParticipantRow, error) {
	results := []fetchParticipantRow{}
	clubs, err := u.repository.EventRegistrationFetchClubByEventID(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("error in fetch clubs : %s", err.Error())
	}
	for _, club := range clubs {
		// totalPoint, _ := u.repository.RankFetchPointByClubID(ctx, club.ID) // temporary comment
		totalUser, _ := u.repository.EventRegistrationCountByEventIDAndClubID(ctx, db.EventRegistrationCountByEventIDAndClubIDParams{
			EventID: eventID,
			ClubID:  club.ID,
		})
		members, err := u.repository.EventParticipantFetchByEventIDAndClubID(ctx, db.EventParticipantFetchByEventIDAndClubIDParams{
			EventID: eventID,
			ClubID:  club.ID,
		})
		if err != nil {
			return nil, fmt.Errorf("error in fetch members : %s", err.Error())
		}
		results = append(results, fetchParticipantRow{
			EventRegistrationFetchClubByEventIDRow: club,
			// TotalPoint:                             totalPoint, // temporary comment
			TotalUser: totalUser,
			Members:   members,
		})
	}

	return results, nil
}
