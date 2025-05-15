package score

import (
	"context"
	"fmt"
	"net/http"

	"github.com/project-ippl-dev/tanding-api/internal/db"
)

type Usecase struct {
	repository    *db.Queries
	rawRepository RawRepository
}

func NewUsecase(repository *db.Queries, rawRepository RawRepository) Usecase {
	return Usecase{repository: repository, rawRepository: rawRepository}
}

func (u Usecase) storeOrUpdateOrder(ctx context.Context, arg orderStoreOrUpdateParams) (statusCode int, err error) {
	orderBracket, err := u.repository.OrderBracketCheckOne(ctx, arg.OrderBracketID)
	if err != nil {
		return http.StatusNotFound, fmt.Errorf("order bracket not found : %s", err.Error())
	}
	if orderBracket.ScoreLock {
		return http.StatusBadRequest, fmt.Errorf("can't update or store score if class event score is already lock")
	}
	if _, err := u.repository.OrderScoreCheckOneByBracketID(ctx, orderBracket.ID); err != nil {
		if err := u.repository.OrderScoreCreate(ctx, db.OrderScoreCreateParams{
			OrderBracketID: orderBracket.ID,
			Round1:         arg.Round1,
			Round2:         arg.Round2,
			Round3:         arg.Round3,
			Extra:          arg.Extra,
			Total:          arg.Total,
		}); err != nil {
			return http.StatusInternalServerError, fmt.Errorf("error in store order score : %s", err.Error())
		}
	} else {
		if err := u.repository.OrderScoreUpdate(ctx, db.OrderScoreUpdateParams{
			Round1:         arg.Round1,
			Round2:         arg.Round2,
			Round3:         arg.Round3,
			Extra:          arg.Extra,
			Total:          arg.Total,
			OrderBracketID: arg.OrderBracketID,
		}); err != nil {
			return http.StatusInternalServerError, fmt.Errorf("error in update order score : %s", err.Error())
		}
	}
	return http.StatusOK, nil
}

func (u Usecase) fetchOneOrder(ctx context.Context, arg fetchOneParams) (db.OrderScoreFetchOneByBracketIDRow, error) {
	return u.repository.OrderScoreFetchOneByBracketID(ctx, arg.BracketID)
}

func (u Usecase) lock(ctx context.Context, arg lockParams) (statusCode int, err error) {
	classEvent, err := u.repository.ClassEventFetchOne(ctx, arg.ClassEventID)
	if err != nil {
		return http.StatusNotFound, fmt.Errorf("class event not found : %s", err.Error())
	}
	if arg.Status == nil {
		*arg.Status = false
	}
	if classEvent.ScoreLock == *arg.Status {
		return http.StatusBadRequest, fmt.Errorf("class event score lock status already %t", *arg.Status)
	}
	if err := u.repository.ClassEventUpdateScoreLock(ctx, db.ClassEventUpdateScoreLockParams{
		ScoreLock: *arg.Status,
		ID:        classEvent.ID,
	}); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error in update score lock : %s", err.Error())
	}
	return http.StatusOK, nil
}

func (u Usecase) storeOrUpdateSingle(ctx context.Context, arg singleStoreOrUpdateParams) (statusCode int, err error) {
	eventBracket, err := u.repository.EventBracketCheckOne(ctx, arg.EventBracketID)
	if err != nil {
		return http.StatusNotFound, fmt.Errorf("single bracket not found : %s", err.Error())
	}
	if eventBracket.ScoreLock {
		return http.StatusBadRequest, fmt.Errorf("can't update or store score if class event score is already lock")
	}

	if eventBracket.IsActive == 0 {
		return http.StatusForbidden, fmt.Errorf("can't store or update event bracket if bracket is not active")
	}

	tx, err := u.rawRepository.db.Begin()
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error in start tx : %s", err.Error())
	}
	txQuery := u.repository.WithTx(tx)
	if _, err := txQuery.EventScoreCheckOneByBracketID(ctx, eventBracket.ID); err != nil {
		if err := txQuery.EventScoreCreate(ctx, db.EventScoreCreateParams{
			EventBracketID: eventBracket.ID,
			HomeRound1:     arg.HomeRound1,
			HomeRound2:     arg.HomeRound2,
			HomeRound3:     arg.HomeRound3,
			HomeExtra:      arg.HomeExtra,
			HomeTotal:      arg.HomeTotal,
			AwayRound1:     arg.AwayRound1,
			AwayRound2:     arg.AwayRound2,
			AwayRound3:     arg.AwayRound3,
			AwayExtra:      arg.AwayExtra,
			AwayTotal:      arg.AwayTotal,
		}); err != nil {
			if err := tx.Rollback(); err != nil {
				return http.StatusInternalServerError, fmt.Errorf("error in rollback tx store event score : %s", err.Error())
			}
			return http.StatusInternalServerError, fmt.Errorf("error in store event score : %s", err.Error())
		}
	} else {
		if err := txQuery.EventScoreUpdate(ctx, db.EventScoreUpdateParams{
			HomeRound1:     arg.HomeRound1,
			HomeRound2:     arg.HomeRound2,
			HomeRound3:     arg.HomeRound3,
			HomeExtra:      arg.HomeExtra,
			HomeTotal:      arg.HomeTotal,
			AwayRound1:     arg.AwayRound1,
			AwayRound2:     arg.AwayRound2,
			AwayRound3:     arg.AwayRound3,
			AwayExtra:      arg.AwayExtra,
			AwayTotal:      arg.AwayTotal,
			EventBracketID: eventBracket.ID,
		}); err != nil {
			if err := tx.Rollback(); err != nil {
				return http.StatusInternalServerError, fmt.Errorf("error in tx rollback update event score : %s", err.Error())
			}
			return http.StatusInternalServerError, fmt.Errorf("error in update event score : %s", err.Error())
		}
	}

	if eventBracket.NextMatchID.String() != "00000000-0000-0000-0000-000000000000" {
		participants, err := txQuery.BracketParticipantFetchByEventBracketID(ctx, eventBracket.ID)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return http.StatusInternalServerError, fmt.Errorf("error in rollback tx bracket participants not found : %s", err.Error())
			}
			return http.StatusInternalServerError, fmt.Errorf("bracket participants not found : %s", err.Error())
		}

		var participantIndex int
		if arg.HomeTotal > arg.AwayTotal {
			participantIndex = 0
		} else {
			participantIndex = 1
		}

		var nextParticipantType db.ParticipantType
		if eventBracket.MatchOrder%2 == 1 {
			nextParticipantType = db.ParticipantTypeHome
		} else {
			nextParticipantType = db.ParticipantTypeAway
		}

		if err := txQuery.BracketParticipantUpdateByParticipantType(ctx, db.BracketParticipantUpdateByParticipantTypeParams{
			EventRegistrationID: participants[participantIndex].EventRegistrationID,
			EventBracketID:      eventBracket.NextMatchID,
			Type:                nextParticipantType,
		}); err != nil {
			if err := tx.Rollback(); err != nil {
				return http.StatusInternalServerError, fmt.Errorf("error in rollback tx update bracket participants update : %s", err.Error())
			}
			return http.StatusInternalServerError, fmt.Errorf("error in update bracket participants : %s", err.Error())
		}

		//Next Bracket
		if err := txQuery.EventBracketUpdateStatus(ctx, db.EventBracketUpdateStatusParams{
			Status:   db.BracketTypeBattle,
			IsActive: 1,
			ID:       eventBracket.NextMatchID,
		}); err != nil {
			if err := tx.Rollback(); err != nil {
				return http.StatusInternalServerError, fmt.Errorf("error in rollback tx update bracket status : %s", err.Error())
			}
			return http.StatusInternalServerError, fmt.Errorf("error in update bracket status : %s", err.Error())
		}
	}

	if err := tx.Commit(); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error in commit tx : %s", err.Error())
	}

	return http.StatusOK, nil
}

func (u Usecase) fetchOneSingle(ctx context.Context, arg fetchOneParams) (db.EventScoreFetchOneByBracketIDRow, error) {
	return u.repository.EventScoreFetchOneByBracketID(ctx, arg.BracketID)
}
