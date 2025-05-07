package event

import (
	"context"
	"database/sql"
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
	return Usecase{repository: repository, rawRepository: rawRepository}
}

func (u Usecase) store(ctx context.Context, req request, decoded tools.JWT) (statusCode int, eventID uuid.UUID, err error) {
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return http.StatusBadRequest, uuid.UUID{}, fmt.Errorf("error in parsing time : " + err.Error())
	}
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return http.StatusBadRequest, uuid.UUID{}, fmt.Errorf("error in parsing time : " + err.Error())
	}
	deadline, err := time.Parse("2006-01-02T15:04:05", req.Deadline)
	if err != nil {
		return http.StatusBadRequest, uuid.UUID{}, fmt.Errorf("error in parsing time : " + err.Error())
	}
	open, err := time.Parse("2006-01-02T15:04:05", req.Open)
	if err != nil {
		return http.StatusBadRequest, uuid.UUID{}, fmt.Errorf("error in parsing time : " + err.Error())
	}

	if deadline.Before(open) {
		return http.StatusUnprocessableEntity, uuid.UUID{}, fmt.Errorf("open must before deadline time")
	}

	tx, err := u.rawRepository.db.Begin()
	if err != nil {
		return http.StatusInternalServerError, uuid.UUID{}, fmt.Errorf("error in start transaction : " + err.Error())
	}
	txQuery := u.repository.WithTx(tx)
	eventID, err = txQuery.EventCreate(ctx, db.EventCreateParams{
		UserID:       uuid.MustParse(decoded.ID),
		Type:         req.Type,
		Name:         req.Name,
		Description:  req.Description,
		PrizePool:    req.PrizePool,
		Location:     req.Location,
		Province:     req.Province,
		City:         req.City,
		Thumbnail:    req.Thumbnail,
		StartDate:    startDate,
		EndDate:      endDate,
		Deadline:     deadline,
		SportID:      uuid.MustParse(req.SportID),
		Rules:        req.Rules,
		ProposalLink: req.ProposalLink,
		Quota:        req.Quota,
		Open:         open,
	})
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return http.StatusInternalServerError, uuid.UUID{}, fmt.Errorf("error in rollback tx : " + err.Error())
		}
		return http.StatusInternalServerError, uuid.UUID{}, fmt.Errorf("error in create event : " + err.Error())
	}

	if err := txQuery.EventPrivilegeCreate(ctx, db.EventPrivilegeCreateParams{
		EventID: eventID,
		UserID:  uuid.MustParse(decoded.ID),
		Role:    db.EventRoleOwner,
	}); err != nil {
		if err := tx.Rollback(); err != nil {
			return http.StatusInternalServerError, uuid.UUID{}, fmt.Errorf("error in rollback tx : " + err.Error())
		}
		return http.StatusInternalServerError, uuid.UUID{}, fmt.Errorf("error in create event privilege : " + err.Error())
	}

	if err := tx.Commit(); err != nil {
		return http.StatusInternalServerError, uuid.UUID{}, fmt.Errorf("error in commit tx : " + err.Error())
	}

	return http.StatusCreated, eventID, nil
}

func (u Usecase) fetchAll(ctx context.Context, page int32, pageSize int32) (tools.Pagination, error) {
	skip := tools.PaginationSkip(page, pageSize)
	events, err := u.repository.EventFetchAll(ctx, db.EventFetchAllParams{
		Limit:  pageSize,
		Offset: skip,
	})
	if err != nil {
		return tools.Pagination{}, fmt.Errorf("error in fetch all event : " + err.Error())
	}
	count, err := u.repository.EventCountAll(ctx)
	if err != nil {
		return tools.Pagination{}, fmt.Errorf("error in count all event : " + err.Error())
	}

	data := []response{}
	for _, event := range events {
		classEvents, err := u.repository.ClassEventFetchByEventID(ctx, event.ID)
		if err != nil {
			return tools.Pagination{}, fmt.Errorf("error in fetch class event : " + err.Error())
		}
		data = append(data, response{
			ID:           event.ID,
			UserID:       event.UserID,
			UserName:     event.UserName,
			UserImage:    event.UserImage,
			Type:         event.Type,
			Name:         event.Name,
			Description:  event.Description,
			PrizePool:    event.PrizePool,
			Location:     event.Location,
			Province:     event.Province,
			City:         event.City,
			Thumbnail:    event.Thumbnail,
			StartDate:    event.StartDate.Format("Monday, 02 January 2006"),
			EndDate:      event.EndDate.Format("Monday, 02 January 2006"),
			Deadline:     event.Deadline.Format("02 January 2006, 15:04"),
			SportID:      event.SportID,
			SportName:    event.SportName,
			Rules:        event.Rules,
			ProposalLink: event.ProposalLink,
			Status:       event.Status.Bool,
			Quota:        event.Quota,
			Open:         event.Open.Format("Monday, 02 January 2006, 15:04"),
			Remark:       string(event.Remark),
			ClassEvents:  classEvents,
		})
	}

	return tools.Pagination{
		TotalItem: count,
		PageSize:  pageSize,
		Page:      page,
		Data:      data,
	}, nil
}

func (u Usecase) update(ctx context.Context, req request, eventID uuid.UUID) (statusCode int, err error) {
	event, err := u.repository.EventFetchOne(ctx, eventID)
	if err != nil {
		return http.StatusNotFound, fmt.Errorf("error in check event : " + err.Error())
	}
	if event.Remark == db.RemarkTypeDone {
		return http.StatusForbidden, fmt.Errorf("can't update event when remark status is done")
	}
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("error in parsing time : " + err.Error())
	}
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("error in parsing time : " + err.Error())
	}
	deadline, err := time.Parse("2006-01-02T15:04:05", req.Deadline)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("error in parsing time : " + err.Error())
	}
	open, err := time.Parse("2006-01-02T15:04:05", req.Open)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("error in parsing time : " + err.Error())
	}

	if deadline.Before(open) {
		return http.StatusUnprocessableEntity, fmt.Errorf("open must before deadline time")
	}

	tx, err := u.rawRepository.db.Begin()
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error in start transaction : " + err.Error())
	}

	txQuery := u.repository.WithTx(tx)
	if err := txQuery.EventUpdate(ctx, db.EventUpdateParams{
		Name:         req.Name,
		Type:         req.Type,
		Description:  req.Description,
		PrizePool:    req.PrizePool,
		Location:     req.Location,
		Province:     req.Province,
		City:         req.City,
		Thumbnail:    req.Thumbnail,
		StartDate:    startDate,
		EndDate:      endDate,
		Deadline:     deadline,
		SportID:      uuid.MustParse(req.SportID),
		Rules:        req.Rules,
		ProposalLink: req.ProposalLink,
		Quota:        req.Quota,
		Open:         open,
		ID:           eventID,
	}); err != nil {
		if err := tx.Rollback(); err != nil {
			return http.StatusInternalServerError, fmt.Errorf("error in rollback tx : " + err.Error())
		}
		return http.StatusInternalServerError, fmt.Errorf("error in update event : " + err.Error())
	}

	now := time.Now()
	if open.Before(now) && event.Status.Bool && event.Remark == "soon" {
		if err := txQuery.EventUpdateRemark(ctx, db.EventUpdateRemarkParams{
			Remark: "open",
			ID:     eventID,
		}); err != nil {
			if err := tx.Rollback(); err != nil {
				return http.StatusInternalServerError, fmt.Errorf("error in rollback tx : " + err.Error())
			}
			return http.StatusInternalServerError, fmt.Errorf("error in update remark event : " + err.Error())
		}
	}

	if err := tx.Commit(); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error in commit tx : " + err.Error())
	}

	return http.StatusOK, nil
}

func (u Usecase) delete(ctx context.Context, eventID uuid.UUID) error {
	if _, err := u.repository.EventCheckOne(ctx, eventID); err != nil {
		return fmt.Errorf("error in check event : " + err.Error())
	}
	return u.repository.EventDelete(ctx, eventID)
}

func (u Usecase) fetchByUser(ctx context.Context, page int32, pageSize int32, userID string) (tools.Pagination, error) {
	skip := tools.PaginationSkip(page, pageSize)
	events, err := u.repository.EventFetchByUserID(ctx, db.EventFetchByUserIDParams{
		UserID: uuid.MustParse(userID),
		Limit:  pageSize,
		Offset: skip,
	})
	if err != nil {
		return tools.Pagination{}, fmt.Errorf("error in fetch event by user id : " + err.Error())
	}
	count, err := u.repository.EventCountByUserID(ctx, uuid.MustParse(userID))
	if err != nil {
		return tools.Pagination{}, fmt.Errorf("error in count event by user id : " + err.Error())
	}

	data := []response{}
	for _, event := range events {
		// totalParticipants, err := u.repository.EventRegistrationCountAllByStatusApproved(ctx, event.ID)  // temporary comment
		// if err != nil {
		// 	return tools.Pagination{}, fmt.Errorf("error in count all participants : " + err.Error())
		// }
		data = append(data, response{
			ID:           event.ID,
			UserID:       event.UserID,
			UserName:     event.UserName,
			UserImage:    event.UserImage,
			Type:         event.Type,
			Name:         event.Name,
			Description:  event.Description,
			PrizePool:    event.PrizePool,
			Location:     event.Location,
			Province:     event.Province,
			City:         event.City,
			Thumbnail:    event.Thumbnail,
			StartDate:    event.StartDate.Format("Monday, 02 January 2006"),
			EndDate:      event.EndDate.Format("Monday, 02 January 2006"),
			Deadline:     event.Deadline.Format("02 January 2006, 15:04"),
			SportID:      event.SportID,
			SportName:    event.SportName,
			Rules:        event.Rules,
			ProposalLink: event.ProposalLink,
			Status:       event.Status.Bool,
			Quota:        event.Quota,
			Open:         event.Open.Format("Monday, 02 January 2006, 15:04"),
			Remark:       string(event.Remark),
			ClassEvents:  nil,
			// Participants: totalParticipants, // temporary comment
		})
	}

	return tools.Pagination{
		TotalItem: count,
		PageSize:  pageSize,
		Page:      page,
		Data:      data,
	}, nil
}

func (u Usecase) fetchOne(ctx context.Context, eventID uuid.UUID, userID string) (responseFetchOne, error) {
	event, err := u.repository.EventFetchOne(ctx, eventID)
	if err != nil {
		return responseFetchOne{}, fmt.Errorf("error in fetch one : " + err.Error())
	}

	classes, err := u.repository.ClassEventFetchByEventID(ctx, event.ID)
	if err != nil {
		return responseFetchOne{}, fmt.Errorf("error in fetch class events : " + err.Error())
	}

	classEvents := []classEventSummaryRow{}
	for _, class := range classes {
		//ranks, err := u.rawRepository.RankFetchByClassEventID(ctx, class.ID) // temporary comment
		// if err != nil {
		// 	return responseFetchOne{}, fmt.Errorf("error in fetch ranks : " + err.Error())
		// }
		classEvents = append(classEvents, classEventSummaryRow{
			ClassEventFetchByEventIDRow: class,
			//Summary:                     ranks,
		})
	}

	// clubRanks, err := u.repository.RankClubFetchByEventID(ctx, eventID)
	// if err != nil {
	// 	return responseFetchOne{}, fmt.Errorf("error fetch in club ranks : " + err.Error())
	// } // temporary comment

	// generalChampions := []db.RankClubFetchByEventIDRow{}
	// for i, clubRank := range clubRanks {
	// 	if i > 2 {
	// 		break
	// 	}
	// 	generalChampions = append(generalChampions, clubRank)
	// } // temporary comment

	// totalParticipants, err := u.repository.EventRegistrationCountAllByStatusApproved(ctx, event.ID)
	// if err != nil {
	// 	return responseFetchOne{}, fmt.Errorf("error in count all participants : " + err.Error())
	// } // temporary comment

	privilege, _ := u.repository.EventPrivilegeFetchOne(ctx, db.EventPrivilegeFetchOneParams{
		EventID: event.ID,
		UserID:  uuid.MustParse(userID),
	})

	return responseFetchOne{
		ID:           event.ID,
		UserID:       event.UserID,
		UserName:     event.UserName,
		Type:         event.Type,
		Name:         event.Name,
		Description:  event.Description,
		PrizePool:    event.PrizePool,
		Location:     event.Location,
		Province:     event.Province,
		City:         event.City,
		Thumbnail:    event.Thumbnail,
		StartDate:    event.StartDate.Format("Monday, 02 January 2006"),
		EndDate:      event.EndDate.Format("Monday, 02 January 2006"),
		Deadline:     event.Deadline,
		SportID:      event.SportID,
		SportName:    event.SportName,
		Rules:        event.Rules,
		ProposalLink: event.ProposalLink,
		Status:       event.Status.Bool,
		Quota:        event.Quota,
		Open:         event.Open.Format("Monday, 02 January 2006, 15:04"),
		Remark:       string(event.Remark),
		ClassEvents:  classEvents,
		UserImage:    event.UserImage,
		// Participants:     totalParticipants, // temporary comment
		Privilege:     privilege,
		EventTurnLock: event.IsGenerate,
		// GeneralChampions: generalChampions, // temporary comment
	}, nil
}

func (u Usecase) fetchInfinite(ctx context.Context, limit int32, args fetchInfiniteQueryParams) (responseInfinite, error) {
	if args.Order == 0 {
		latestOrder, _ := u.rawRepository.EventFetchLatestOrder(ctx, fetchQueryParams{
			SportID:  args.SportID,
			Name:     "%" + args.Name + "%",
			Category: args.Category,
			Remark:   args.Remark,
		})
		args.Order = latestOrder + 1
	}
	events, err := u.rawRepository.EventFetchInfinite(ctx, fetchInfiniteParams{
		OrderNumber: args.Order,
		Limit:       limit,
		fetchQueryParams: fetchQueryParams{
			SportID:  args.SportID,
			Name:     "%" + args.Name + "%",
			Category: args.Category,
			Remark:   args.Remark,
		},
	})
	if err != nil {
		return responseInfinite{}, fmt.Errorf("error in fetch events : " + err.Error())
	}

	total, err := u.repository.EventCountInfinite(ctx)
	if err != nil {
		return responseInfinite{}, fmt.Errorf("error in count infinite : " + err.Error())
	}

	data := []responseInfiniteRow{}
	for _, event := range events {
		// totalParticipants, err := u.repository.EventRegistrationCountAllByStatusApproved(ctx, event.ID) // temporary comment
		// if err != nil {
		// 	return responseInfinite{}, fmt.Errorf("error in count all participants : " + err.Error())
		// }
		data = append(data, responseInfiniteRow{
			ID:          event.ID,
			UserID:      event.UserID,
			UserName:    event.UserName,
			Type:        event.Type,
			Name:        event.Name,
			Description: event.Description,
			PrizePool:   event.PrizePool,
			Location:    event.Location,
			Province:    event.Province,
			City:        event.City,
			Thumbnail:   event.Thumbnail,
			StartDate:   event.StartDate.Format("Monday, 02 January 2006"),
			EndDate:     event.EndDate.Format("Monday, 02 January 2006"),
			Deadline:    event.Deadline.Format("02 January 2006, 15:04"),
			SportID:     event.SportID,
			SportName:   event.SportName,
			Quota:       event.Quota,
			Order:       event.OrderNumber,
			Open:        event.Open.Format("Monday, 02 January 2006, 15:04"),
			Remark:      string(event.Remark),
			// Participants: totalParticipants, // temporary comment
			UserImage: event.UserImage,
		})
	}
	return responseInfinite{
		Message:   "fetch infinite scroll for event success",
		Data:      data,
		TotalItem: total,
	}, nil
}

func (u Usecase) updateStatus(ctx context.Context, req statusReq, eventID uuid.UUID) error {
	event, err := u.repository.EventFetchOne(ctx, eventID)
	if err != nil {
		return fmt.Errorf("error in fetch one event :" + err.Error())
	}

	var remark db.RemarkType
	if !event.Status.Bool {
		if event.Open.Before(time.Now()) {
			remark = db.RemarkTypeOpen
		} else {
			remark = db.RemarkTypeSoon
		}
	} else {
		remark = event.Remark
	}
	if *req.Status {
		if err := u.repository.EventUpdateStatus(ctx, db.EventUpdateStatusParams{
			Status: sql.NullBool{
				Bool:  *req.Status,
				Valid: true,
			},
			Remark: remark,
			ID:     eventID,
		}); err != nil {
			return fmt.Errorf("error in update status :" + err.Error())
		}
	} else {
		if err := u.repository.EventUpdateStatus(ctx, db.EventUpdateStatusParams{
			Status: sql.NullBool{
				Bool:  *req.Status,
				Valid: true,
			},
			Remark: db.RemarkTypeRejected,
			ID:     eventID,
		}); err != nil {
			return fmt.Errorf("error in update status :" + err.Error())
		}
	}
	return nil
}

func (u Usecase) assign(ctx context.Context, req assignRequest) error {
	tx, err := u.rawRepository.db.Begin()
	if err != nil {
		return fmt.Errorf("error in start tx : " + err.Error())
	}
	txQuery := u.repository.WithTx(tx)
	for _, data := range req.Data {
		if err := txQuery.EventPrivilegeCreate(ctx, db.EventPrivilegeCreateParams{
			EventID: uuid.MustParse(req.EventID),
			UserID:  uuid.MustParse(data.UserID),
			Role:    data.Role,
		}); err != nil {
			if err := tx.Rollback(); err != nil {
				return fmt.Errorf("error in rollback tx : " + err.Error())
			}
			return fmt.Errorf("error in store event privilege : " + err.Error())
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error in commit tx : " + err.Error())
	}
	return nil
}

func (u Usecase) committeeFetchAll(ctx context.Context, eventID uuid.UUID) ([]db.EventPrivilegeFetchAllRow, error) {
	return u.repository.EventPrivilegeFetchAll(ctx, eventID)
}

func (u Usecase) committeeUpdate(ctx context.Context, arg updateCommitteeParams, userID string) (statusCode int, err error) {
	authPrivilege, err := u.repository.EventPrivilegeFetchOne(ctx, db.EventPrivilegeFetchOneParams{
		EventID: arg.EventID,
		UserID:  uuid.MustParse(userID),
	})
	if _, err := u.repository.EventPrivilegeFetchOne(ctx, db.EventPrivilegeFetchOneParams{
		EventID: arg.EventID,
		UserID:  arg.CommitteeID,
	}); err != nil {
		return http.StatusNotFound, fmt.Errorf("privilege not found : " + err.Error())
	}

	switch authPrivilege.Role {
	case db.EventRoleAdmin:
		if arg.Role == db.EventRoleAdmin || arg.Role == db.EventRoleOwner {
			return http.StatusForbidden, fmt.Errorf("role admin only can update privilege to contributor and reviewer")
		}
		break
	}

	if err := u.repository.EventPrivilegeUpdate(ctx, db.EventPrivilegeUpdateParams{
		Role:    arg.Role,
		UserID:  arg.CommitteeID,
		EventID: arg.EventID,
	}); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error in update event privilege : " + err.Error())
	}

	return http.StatusOK, nil
}

func (u Usecase) committeeDelete(ctx context.Context, arg deleteCommitteeParams, userID string) (statusCode int, err error) {
	authPrivilege, err := u.repository.EventPrivilegeFetchOne(ctx, db.EventPrivilegeFetchOneParams{
		EventID: arg.EventID,
		UserID:  uuid.MustParse(userID),
	})
	privilege, err := u.repository.EventPrivilegeFetchOne(ctx, db.EventPrivilegeFetchOneParams{
		EventID: arg.EventID,
		UserID:  arg.CommitteeID,
	})
	if err != nil {
		return http.StatusNotFound, fmt.Errorf("privilege not found : " + err.Error())
	}

	switch authPrivilege.Role {
	case db.EventRoleAdmin:
		if privilege.Role == db.EventRoleAdmin || privilege.Role == db.EventRoleOwner {
			return http.StatusForbidden, fmt.Errorf("role admin only can delete privilege to contributor and reviewer")
		}
		break
	}

	if err := u.repository.EventPrivilegeDelete(ctx, db.EventPrivilegeDeleteParams{
		UserID:  arg.CommitteeID,
		EventID: arg.EventID,
	}); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error in delete event privilege : " + err.Error())
	}

	return http.StatusOK, nil
}

func (u Usecase) updateRemark(ctx context.Context, arg updateRemarkParams) error {
	event, err := u.repository.EventCheckOne(ctx, arg.EventID)
	if err != nil {
		return fmt.Errorf("event not found : " + err.Error())
	}
	if err := u.repository.EventUpdateRemark(ctx, db.EventUpdateRemarkParams{
		Remark: arg.Remark,
		ID:     event.ID,
	}); err != nil {
		return fmt.Errorf("error in update remark event : " + err.Error())
	}
	return nil
}
