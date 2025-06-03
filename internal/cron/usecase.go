package cron

import (
	"context"
	"fmt"
	"time"

	"github.com/project-ippl-dev/tanding-api/internal/db"
)

type Usecase struct {
	repository    *db.Queries
	rawRepository RawRepository
}

func NewUsecase(repository *db.Queries, rawRepository RawRepository) Usecase {
	return Usecase{repository: repository, rawRepository: rawRepository}
}

func (u Usecase) eventUpdateRemarkSoonToOpen() error {
	now := time.Now()
	ctx := context.Background()
	tx, err := u.rawRepository.db.Begin()
	if err != nil {
		return fmt.Errorf("error in start transaction : %s", err.Error())
	}
	txQuery := u.repository.WithTx(tx)
	events, err := txQuery.EventFetchByRemarkSoon(ctx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf("error in rollback tx : %s", err.Error())
		}
		return fmt.Errorf("error in fetch event by remark soon : %s", err.Error())
	}
	for _, event := range events {
		if event.Open.Unix() <= now.Unix() {
			if err := txQuery.EventUpdateRemark(ctx, db.EventUpdateRemarkParams{
				Remark: db.RemarkTypeOpen,
				ID:     event.ID,
			}); err != nil {
				if err := tx.Rollback(); err != nil {
					return fmt.Errorf("error in rollback tx event update remark soon to open : %s", err.Error())
				}
				return fmt.Errorf("error in update remark event soon to open : %s", err.Error())
			}
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error in commit tx : %s", err.Error())
	}

	return nil
}

func (u Usecase) eventUpdateRemarkOpenToClose() error {
	now := time.Now()
	ctx := context.Background()
	tx, err := u.rawRepository.db.Begin()
	if err != nil {
		return fmt.Errorf("error in start transaction : %s", err.Error())
	}
	txQuery := u.repository.WithTx(tx)
	events, err := txQuery.EventFetchByRemarkOpen(ctx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf("error in rollback tx : %s", err.Error())
		}
		return fmt.Errorf("error in fetch event by remark open to closed : %s", err.Error())
	}
	for _, event := range events {
		if event.Deadline.Unix() <= now.Unix() {
			if err := txQuery.EventUpdateRemark(ctx, db.EventUpdateRemarkParams{
				Remark: db.RemarkTypeClosed,
				ID:     event.ID,
			}); err != nil {
				if err := tx.Rollback(); err != nil {
					return fmt.Errorf("error in rollback tx event update remark open to closed : %s", err.Error())
				}
				return fmt.Errorf("error in update remark event open to closed : %s", err.Error())
			}
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error in commit tx : %s", err.Error())
	}

	return nil
}

func (u Usecase) eventUpdateRemarkCloseToOngoing() error {
	now := time.Now()
	ctx := context.Background()
	tx, err := u.rawRepository.db.Begin()
	if err != nil {
		return fmt.Errorf("error in start transaction : %s", err.Error())
	}
	txQuery := u.repository.WithTx(tx)
	events, err := txQuery.EventFetchByRemarkClose(ctx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf("error in rollback tx : %s", err.Error())
		}
		return fmt.Errorf("error in fetch event by remark closed to ongoing : %s", err.Error())
	}
	for _, event := range events {
		if event.StartDate.Unix() <= now.Unix() {
			if err := txQuery.EventUpdateRemark(ctx, db.EventUpdateRemarkParams{
				Remark: db.RemarkTypeOngoing,
				ID:     event.ID,
			}); err != nil {
				if err := tx.Rollback(); err != nil {
					return fmt.Errorf("error in rollback tx event update remark closed to ongoing : %s", err.Error())
				}
				return fmt.Errorf("error in update remark event closed to ongoing : %s", err.Error())
			}
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error in commit tx : %s", err.Error())
	}

	return nil
}

//func (u Usecase) eventUpdateRemarkOngoingToDone() error {
//	now := time.Now()
//	ctx := context.Background()
//	tx, err := u.rawRepository.db.Begin()
//	if err != nil {
//		return fmt.Errorf("error in start transaction : %s", err.Error())
//	}
//	txQuery := u.repository.WithTx(tx)
//	events, err := txQuery.EventFetchByRemarkOngoing(ctx)
//	if err != nil {
//		if err := tx.Rollback(); err != nil {
//			return fmt.Errorf("error in rollback tx : %s", err.Error())
//		}
//		return fmt.Errorf("error in fetch event by remark ongoing to done : %s", err.Error())
//	}
//	for _, event := range events {
//		if event.EndDate.Unix() <= now.Unix() {
//			if err := txQuery.EventUpdateRemark(ctx, db.EventUpdateRemarkParams{
//				Remark: db.RemarkTypeDone,
//				ID:     event.ID,
//			}); err != nil {
//				if err := tx.Rollback(); err != nil {
//					return fmt.Errorf("error in rollback tx event update remark ongoing to done : %s", err.Error())
//				}
//				return fmt.Errorf("error in update remark event ongoing to done : %s", err.Error())
//			}
//		}
//	}
//	if err := tx.Commit(); err != nil {
//		return fmt.Errorf("error in commit tx : %s", err.Error())
//	}
//
//	return nil
//}

func (u Usecase) registrationUpdate() error {
	now := time.Now()

	registrations, err := u.repository.EventRegistrationFetchPendingCron(context.Background(), now)
	if err != nil {
		return fmt.Errorf("error in fetch registrations : %s", err.Error())
	}
	tx, err := u.rawRepository.db.Begin()
	if err != nil {
		return fmt.Errorf("error in start tx registration update : %s", err.Error())
	}
	txQuery := u.repository.WithTx(tx)
	for _, registration := range registrations {
		if err := txQuery.EventRegistrationUpdateStatus(context.Background(), db.EventRegistrationUpdateStatusParams{
			Status: db.EventRegistrationStatusCanceled,
			ID:     registration,
		}); err != nil {
			if err := tx.Rollback(); err != nil {
				return fmt.Errorf("error in rollback tx registration update status : %s", err.Error())
			}
			return fmt.Errorf("error in registration update status : %s", err.Error())
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error in commit tx registration update status : %s", err.Error())
	}

	return nil
}
